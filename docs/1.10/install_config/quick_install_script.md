# Deploy Harbor with the Quick Installation Script

The Harbor community has provided a script that with a single command prepares an Ubuntu 18.04 machine for Harbor and deploys the latest stable version.

## Prerequisites

You have a machine or VM that is running Ubuntu 18.04.

## Procedure

1. Download the `harbor.sh` script from [gist.github.com](https://gist.github.com/kacole2/95e83ac84fec950b1a70b0853d6594dc) to your Ubuntu machine or VM.
1. Grant run permissions to the current user.

   ```$ chmod u+x```
1. Run the script as superuser.

   ```$ sudo ./harbor.sh```
1. Select whether to deploy Harbor using the IP address or FQDN of the host machine. 

   This is the address at which you access the Harbor interface and the registry service.
   
   - To use the IP address, enter `1`.
   - To use the FQDN, enter `2`.
   
   The script takes several minutes to run. As it runs, the script downloads the necessary packages and dependencies from Ubuntu, installs the latest stable versions of Docker and Docker Compose, and installs the latest stable version of Harbor.
   
1. When the script reports `Harbor Installation Complete`, log in to your new Harbor instance. 

   ```$ docker login <harbor_ip_or_FQDN>```
   
   - User name: `admin`
   - Password: `VMware12345` 
1. Enter the Harbor address in a browser to log in to the Harbor interface.