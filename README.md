# zitadel-bootstrapper

Setup and configure Zitadel and components using OIDC

# Components

## Admin account

Create a Zitadel admin account

```yaml
config:
  adminAccount:
    setup: false
    firstname: Platform
    lastname: Person
    username: platform@fabled.se
    password: ""
```

## ArgoCD

Creates a new Zitadel project and OIDC application. The OIDC `clientId` and `clientSecret` will be inside a secret `argocd-zitadel-oidc` within the `argocd` namespace.

```yaml
config:
  argoCD:
    setup: false
    name: ArgoCD
    userRoleName: argocd_users
    adminRoleName: argocd_administrators
    devMode: false # Needs to be true if plain http redirect/logout uri is used
    redirectUris:
      - https://argocd.argocd.svc.cluster.local/auth/callback
    logoutUris:
      - https://argocd.argocd.svc.cluster.local
```

Example usage of argo-cd helm chart
```yaml
      configs:
        cm:
          oidc.config: |
            name: Zitadel
            issuer: https://argocd.argocd.svc.cluster.local
            ClientID: $argocd-zitadel-oidc:clientId
            ClientSecret: $argocd-zitadel-oidc:secretId
            requestedScopes:
              - openid
              - profile
              - email
              - groups
            logoutURL: https://argocd.argocd.svc.cluster.local
          /.../
        rbac:
          scopes: "[groups]"
          policy.csv: |
            g, argocd_administrators, role:admin
            g, argocd_users, role:readonly
          policy.default: ''
```
