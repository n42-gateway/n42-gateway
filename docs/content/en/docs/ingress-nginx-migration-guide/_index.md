---
title: "ingress-nginx Migration Guide"
linkTitle: "ingress-nginx Migration Guide"
weight: 3
description: >
  Guide for migrating a Service from being fronted by an ingress-nginx Ingress to an n42-gateway one.
---

This is an opinionated zero-migration guide for migrating Ingress objects from being managed by ingress-nginx controllers to n42-gateway.

Not all Kubernetes clusters are configured with the same controllers that act on Ingress objects, and not all migration requirements are the same. Use this guide to inform your own migration plans. The provided example objects are examples only, and are not expected to work in every environment.

## Configure new n42-gateway controller

1. You will need your old ingress-nginx controller, and a new n42-gateway controller that will NOT pick up the same Ingress objects. Configure the n42-gateway controller to use a new IngressClass with the `--ingress-class` flag. It is important that n42-gateway uses an IngressClass that is different than the one used by the ingress-nginx controller. See [Class matter](https://n42-gateway.github.io/docs/configuration/keys/#class-matter)  to understand how n42-gateway uses IngressClasses.

For example, configure n42-gateway to use the startup flag `--ingress-class=migration-demo` and a dedicated IngressClass:

```yaml
---
apiVersion: networking.k8s.io/v1
kind: IngressClass
metadata:
  name: external-haproxy
spec:
  controller: n42-gateway.github.io/controller/external-haproxy
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: n42-gateway-controller
spec:
  template:
    spec:
      containers:
      - name: controller
        image: localhost/n42-gateway:latest
        args:
        - --ingress-class=external-haproxy
        # ... other arguments like --configmap
```

## Test parallel n42-gateway HAProxy Ingress object

1. Create a parallel Ingress object for your Service. It should be using the new IngressClass and specify a different `.spec.rules[*].host` to make sure there are no conflicts with the ingress-nginx controller or any other controllers that may be acting on the Ingress object.
2. Ensure that the parallel Ingress object is configured correctly for HAProxy as opposed to NGINX. See [Configuration keys](https://n42-gateway.github.io/docs/configuration/keys/) for available configurations. Some configuration keys may be backwards compatible with ingress-nginx, but not all.

For example, an ingress-nginx managed Ingress object and an n42-gateway managed one:

```yaml
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-nginx-ingress
  annotations:
    # Optional: Optional nginx-specific configurations go here
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx # Directs traffic to the ingress-nginx controller
  rules:
  - host: ingress.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: backend-service
            port:
              number: 80
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-n42-ingress
  annotations:
    # Optional: Optional n42/HAProxy-specific configurations go here
    haproxy.org/ssl-redirect: "true"
spec:
  ingressClassName: external-haproxy # Directs traffic to the n42-gateway controller
  rules:
  - host: parallel-ingress.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: backend-service
            port:
              number: 80
```

## Execute Migration

> [!WARNING]
> Kubernetes clusters may have many controllers that also interact with Ingress objects. You must know which ones are relevant to your environment in order to configure Ingresses correctly. Popular controllers could be [cert-manager](https://cert-manager.io/) and [external-dns](https://kubernetes-sigs.github.io/external-dns/latest/), but there could be many more. Particularly with DNS controllers, this could complicate the steps to achieve a zero-downtime migration.

1. Configure the new parallel Ingress by updating `.spec.rules[*].host` to match with the old Ingress object.

   If you are not using DNS records such as `ingress.example.com` to access your Kubernetes Service via the Ingress, you can skip to step 4. Otherwise, continue on.

2. Update DNS record value for the Ingress URL to point to the new n42-gateway controller.

   Typically, the DNS record for an Ingress such as `ingress.example.com` will point to the External IP address of the Service corresponding to the Ingress controller.
   By updating the DNS record value to be te External IP address of the n42-gateway Service, clients whose DNS caches expire will begin to send requests to n42-gateway,
   and clients that are still using cached values will continue to be served by ingress-nginx until their caches expire as well.

   A gotcha here for clusters that automate DNS record management for Ingress objects is that the old Ingress controller must still be configured to program the host
   for the Ingress object, but the DNS record controller must know that the DNS record value should be assigned to the n42-gateway controller's IP address. Not the old
   ingress-nginx controller's IP address.

3. Observe traffic patterns to watch request traffic migrate away from ingress-nginx to n42-gateway.

   Switching the value of the DNS record for the Ingress object to point to the new n42-gateway controller will take at least as long as the TTL value of said record.
   For at least the duration of the DNS record TTL, observe traffic slowly get served more and more by n42-gateway as opposed to ingress-nginx until ingress-nginx no
   longer serves any traffic.

4. Delete the old Ingress object.

   Once the ingress-nginx controller is no longer serving requests for the Ingress URL, it is safe to delete.
