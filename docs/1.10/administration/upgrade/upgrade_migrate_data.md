[Back to table of contents](../../index.md)

----------

# Upgrade Harbor and Migrate Data

This guide covers upgrade and migration to version 1.10.0. This guide only covers migration from v1.8.x and later to the current version. If you are upgrading from an earlier version, refer to the migration guide in the `release-1.8.0` branch to upgrade to v1.8.x first, then follow this guide to perform the migration to this version.

When upgrading an existing Harbor instance to a newer version, you might need to migrate the data in your database and the settings in `harbor.cfg`.
Since the migration might alter the database schema and the settings of `harbor.cfg`, you should **always** back up your data before any migration.

**NOTES:**

- Again, you must back up your data before any data migration.
- In version 1.9.0, some containers are started by `non-root`. This does not pose problems if you are upgrading an officially released version of Harbor, but if you have deployed a customized instance of Harbor, you might encounter permission issues.
- In previous releases, user roles took precedence over group roles in a project. In this version, user roles and group roles are combined so that the user has whichever set of permissions is highest. This might cause the roles of certain users to change during upgrade.
- With the introduction of storage and artifact quotas in version 1.9.0, migration from 1.8.x might take a few minutes. This is because the `core` walks through all blobs in the registry and populates the database with information about the layers and artifacts in projects.
- With the introduction of storage and artifact quotas in version 1.9.0, replication between version 1.9.0 and a previous version of Harbor does not work. You must upgrade all Harbor nodes to 1.9.0 if you have configured replication between them.

## Upgrading Harbor and Migrating Data

1. Log in to the Harbor host and, if it is still running, stop and remove the existing Harbor instance.

    ```sh
    cd harbor
    docker-compose down
    ```

1. Back up Harbor's current files so that you can roll back to the current version if necessary.

    ```sh
    mv harbor /my_backup_dir/harbor
    ```

1. Back up the database, which by default is in the directory `/data/database`.

    ```sh
    cp -r /data/database /my_backup_dir/
    ```

1. Get the latest Harbor release package from [https://github.com/goharbor/harbor/releases](https://github.com/goharbor/harbor/releases)
1. Before upgrading Harbor, perform migration. 

    The migration tool is delivered as a docker image. You can pull the image from docker hub. Replace [tag] with the new Harbor version, for example v1.10.0, in the following command:
    
    ```sh
    docker pull goharbor/harbor-migrator:[tag]
    ```

    Alternatively, if you are using an offline installer package, you can load it from the image tarball that is included in the offline installer package. Replace [tag] with the new Harbor version, for example v1.10.0, in the following command:
    
    ```sh
    tar zxf <offline package>
    docker image load -i harbor/harbor.[version].tar.gz
    ```

1. Upgrade the `harbor.yml` file.

    ```sh
    docker run -it --rm -v ${harbor_yml}:/harbor-migration/harbor-cfg/harbor.yml goharbor/harbor-migrator:[tag] --cfg up
    ```

    **NOTE:** The schema upgrade and data migration of the database is performed by core when Harbor starts. If the migration fails, check the core log to debug.

1. In the `./harbor` directory, run the `./install.sh` script to install the new Harbor instance. 

   To install Harbor with components such as Notary, Clair, and chartmuseum, see [Run the Installer Script](../../install_config/run_installer_script.md) for more information.

----------

[Back to table of contents](../../index.md)