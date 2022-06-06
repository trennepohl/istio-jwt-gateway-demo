## Istio JWT validation proof of concept
This project demonstrates how one can login using Google SSO, and validate a self-signed JWT token using Istio across different services

### Requirements
-  [Docker](https://docs.docker.com/get-docker/)
-  [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
-  [Istioctl](https://istio.io/latest/docs/reference/commands/istioctl/)
   -  Arch: `yay -Sy istioctl`
   -  OSX: `brew install istioctl`
-  [Kubectl](https://kubernetes.io/docs/tasks/tools/)
-  [Setup Google Oauth credentials](https://support.google.com/cloud/answer/6158849?hl=en)


### Setting up Google Oauth (Optional, but if not configured SSO with google won't work)
- General guide on how to [setup Google oauth credentials](https://support.google.com/cloud/answer/6158849?hl=en)

The callback url must be `http://authorization.com:8080/auth/callback/google` unless you change the code and manifests to use a different domain or URL. If so please set the URL accordingly.

After the setup is ready, head to the [google auth secret manifest](manifests/authorization-svc.yaml) and replace the secrets at the top of the file.

> Don't forget secrets must be base64 encoded

### Getting the cluster running
To install the cluster and all components run `make install`

 **_NOTE:_**  For the demo to work properly you'll have to add two extra lines to your `hosts` file
 ```
 127.0.0.1 authorization.com
 127.0.0.1 samplesvc.com
 ```

 Once the instalation is complete, running `make proxy` will make the istio-ingressgateway at port 8080.
 
 To verify everything is working you can run the following request:
 ```bash
    curl -XPOST 'http://authorization.com:8080/auth/login' -H'Host: authorization.com' -H'Content-Type: application/json' --data-raw '{"Email": "admin@istio-auth-poc.io","Password": "password"}'
```

You should see a token issued to an Admin user.

### Acessing the cluster
After the cluster is initialized, a `kubeconfig` file will be available the the root folder of this project.
```bash
export KUBECONFIG=kubeconfig
kubectl get pods
```

### Database credentials
The username and password for the postgres database are defined in the [manifests](manifests/authorization-svc.yaml).

### Authorization policies and JWT validation
If you look into the [istio-authorization-policies.yaml](manifests/istio-authorization-policies.yaml) you will notice there are two simple policies, 1 for the `authorization` service and another one for the `samplesvc`

For the `samplesvc` we allow any requests that contain a role which matches `Admin, ReadOnly or ReadWrite`
For the `authorization` service we require that requests made to all endpoints under `/admin` to come from a token that contains the `Admin` role.

This validation is only possible due to the [istio-jwt-validator](manifests/istio-jwt-validator.yaml) which uses the public keys availabe in the endpoint `/jwk` of the authorization service to valide the token signature.

### Endpoints
Below you'll find the API documentation, however I strongly recomend you to use Postman and import [this](postman-collection.json) collection


#### Login

Endpoint for admin only login

**URL** : `/auth/login/`

**Method** : `POST`

**Auth required** : NO

**Data constraints**

```json
{
    "Email": "[valid email address]",
    "password": "[password in plain text]"
}
```

#### Get Users

Endpoint to list all users

**URL** : `/admin/users`

**Method** : `GET`

**Auth required** : Yes

---

#### Get Roles

Endpoint to list all roles

**URL** : `/admin/roles`

**Method** : `GET`

**Auth required** : Yes

---

#### Create Role

Endpoint to list all roles

**URL** : `/admin/role`

**Method** : `POST`

**Auth required** : Yes

**Data constraints**

```json
{
    "Name": "RoleName"
}
```
---
#### Assign Role

Endpoint to assign a role to a user

**URL** : `/admin/role/assign`

**Method** : `POST`

**Auth required** : Yes

**Data constraints**

```json
{
    "UserID": 1,
    "RoleID": 1
}
```
---
#### Remove Role

Endpoint to remove a role from a user

**URL** : `/admin/role/remove`

**Method** : `POST`

**Auth required** : Yes

**Data constraints**

```json
{
    "UserID": 1,
    "RoleID": 1
}
```

### Authorization service environment variables

| Environment      | Description | Default value |
| ----------- | ----------- | ----------- |
| ADMIN_EMAIL      | User to be created at startup       | admin@istio-auth-poc.io |
| ADMIN_PASSWORD   | Password to be assigned to the default admin user        | password       |
| GOOGLE_STATE_CODE   | A unique text value        | somethingunique       |
| GOOGLE_CALLBACK_URL   | The Google Oauth callback url        | authorization.com:8080/auth/google/callback       |
| GOOGLE_CLIENT_ID   | The Oauth client id        |        |
| GOOGLE_CLIENT_SECRET   | The Oauth client secret        |        |
| DATABASE_USER   | The database user credential        | istio-poc       |
| DATABASE_HOST   | The database host addr        | localhost       |
| DATABASE_PORT   | The database port        | 5432       |
| DATABASE_PASSWORD   | The database password        | mysecretpassword       |
| DATABASE_NAME   | The database name        | authorization       |

### To be done

#### Token revocation
Right now this proof of concept doesn't support the revocation of tokens

#### Authorization Policies conditions
Looks like when there are more than one condition in an AuthorizationPolicy rule the evaluation is `OR` based and not `AND` therefore the conditions below will allow traffic to anyone that has a token issued by `istio-auth-poc` regardless of their roles.
```yml
      when:
        - key: request.auth.claims[iss]
          values: ["istio-auth-poc"]
        - key: request.auth.claims[Roles]
          values: ["Admin"]
```