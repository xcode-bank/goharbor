[Back to table of contents](../../_index.md)

----------

# Run the Installer Script

Once you have configured `harbor.yml` and optionally set up a storage backend, you install and start Harbor by using the `install.sh` script. Note that it might take some time for the online installer to download all of the Harbor images from Docker hub.

You can install Harbor in different configurations:

- Just Harbor, without Notary, Clair, or Chart Repository Service
- Harbor with Notary
- Harbor with Clair
- Harbor with Chart Repository Service
- Harbor with two or all three of Notary, Clair, and Chart Repository Service

## Default installation without Notary, Clair, or Chart Repository Service

The default Harbor installation does not include Notary or Clair service. Run the following command

``` 
sudo ./install.sh
```

If the installation succeeds, you can open a browser to visit the Harbor interface at `http://reg.yourdomain.com`, changing `reg.yourdomain.com` to the hostname that you configured in `harbor.yml`. If you did not change them in `harbor.yml`, the default administrator username and password are `admin` and `Harbor12345`.

Log in to the admin portal and create a new project, for example, `myproject`. You can then use Docker commands to log in to Harbor, tag images, and push them to Harbor.

```
docker login reg.yourdomain.com
```
```
docker push reg.yourdomain.com/myproject/myrepo:mytag
```

**IMPORTANT:** 
- If your installation of Harbor uses HTTPS, you must provide the Harbor certificates to the Docker client. For information, see [Configure HTTPS Access to Harbor](configure_https.md#provide-the-certificates-to-harbor-and-docker).
- If your installation of Harbor uses HTTP, you must add the option `--insecure-registry` to your client's Docker daemon and restart the Docker service. For more information, see [Connecting to Harbor via HTTP](#connect_http) below.

## Installation with Notary

To install Harbor with the Notary service, add the `--with-notary` parameter when you run `install.sh`:

```
sudo ./install.sh --with-notary
```

**NOTE**: For installation with Notary, you must configure Harbor to use HTTPS.

For more information about Notary and Docker Content Trust, see [Content Trust](https://docs.docker.com/engine/security/trust/content_trust/) in the Docker documentation.

## Installation with Clair

To install Harbor with Clair service, add the `--with-clair` parameter when you run `install.sh`:

```
sudo ./install.sh --with-clair
```

For more information about Clair, see the [Clair documentation](https://coreos.com/clair/docs/2.0.1/).

By default, Harbor limits the CPU usage of the Clair container to 150000 to avoid it using up all CPU resources. This is defined in the `docker-compose.clair.yml` file. You can modify this file based on your hardware configuration.

## Installation with Chart Repository Service 

To install Harbor with chart repository service, add the `--with-chartmuseum` parameter when you run `install.sh`:

```
sudo ./install.sh --with-chartmuseum
```

## Installation with Notary, Clair, and Chart Repository Service

If you want to install all three of Notary, Clair and chart repository service, specify all of the parameters in the same command:

```
sudo ./install.sh --with-notary --with-clair --with-chartmuseum
```

<a id="connect_http"></a>
## Connecting to Harbor via HTTP

**IMPORTANT:** If your installation of Harbor uses HTTP rather than HTTPS, you must add the option `--insecure-registry` to your client's Docker daemon. By default, the daemon file is located at `/etc/docker/daemon.json`.

For example, add the following to your `daemon.json` file:

<pre>
{
"insecure-registries" : ["<i>myregistrydomain.com</i>:5000", "0.0.0.0"]
}
</pre>

After you update `daemon.json`, you must restart both Docker Engine and Harbor.

1. Restart Docker Engine.

   `systemctl restart docker`
1. Stop Harbor.

   `docker-compose down -v`
1. Restart Harbor.

   `docker-compose up -d`

## What to Do Next ##

- If the installation succeeds, continue to set up Harbor by following the instructions in [Post-Installation Configuration](../configuration/_index.md) and [Initial Configuration in the Harbor UI](../configuration/initial_config_ui.md).
- If you deployed Harbor with HTTP and you want to secure the connections to Harbor, see [Configure HTTPS Access to Harbor](configure_https.md).
- If installation fails, see [Troubleshooting Harbor Installation
](troubleshoot_installation.md).

----------

[Back to table of contents](../../_index.md)