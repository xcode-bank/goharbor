# Harbor Documentation 

This is the main table of contents for the Harbor documentation.

## Harbor Installation and Configuration

This section describes how to install Harbor and perform the required initial configurations. These day 1 operations are performed by the Harbor Administrator.

- [Harbor Installation and Configuration](install_config/_index.md)
  - [Installing Harbor](install_config/installation/_index.md) 
      - [Test Harbor with the Demo Server](install_config/installation/demo_server.md)
      - [Harbor Compatibility List](install_config/installation/harbor_compatibility_list.md)
      - [Harbor Installation Prerequisites](install_config/installation/installation_prereqs.md)
      - **[Configure HTTPS Access to Harbor](install_config/installation/configure_https.md)**
      - [Download the Harbor Installer](install_config/installation/download_installer.md)
      - **[Configure the Harbor YML File](install_config/installation/configure_yml_file.md)**
      - [Run the Installer Script](install_config/installation/run_installer_script.md)
      - **[Troubleshooting Harbor Installation](install_config/installation/troubleshoot_installation.md)**
  - [Post-Installation Configuration](install_config/configuration/_index.md)
      - **[Reconfigure Harbor and Manage the Harbor Lifecycle](install_config/configuration/reconfigure_manage_lifecycle.md)**
      - **[Customize the Harbor Token Service](install_config/configuration/customize_token_service.md)**
      - **[Configure Notary Content Trust](install_config/configuration/configure_notary_content_trust.md)**
  - [Initial Configuration in the Harbor UI](install_config/configuration/initial_config_ui.md)
      - [Configure Authentication](install_config/configuration/configure_authentication.md)
      - **[Administrator Options](install_config/configuration/general_settings.md)**
  
## Harbor Administration

This section describes how to use and maintain Harbor after deployment. These day 2 operations are performed by the Harbor Administrator.

- [Harbor Administration](administration/_index.md)
  - [Manage Users](administration/managing_users/_index.md)
      - **[Harbor Role Based Access Control (RBAC)](administration/managing_users/configure_rbac.md)**
      - [User Permissions By Role](administration/managing_users/user_permissions_by_role.md)
      - **[Configure Harbor User Settings at the Command Line](administration/managing_users/configure_user_settings_cli.md)**
  - [Configure Project Settings](administration/configure_project_settings/_index.md)
      - [Set Project Quotas](administration/configure_project_settings/set_project_quotas.md)
  - [Configuring Replication](administration/configuring_replication/_index.md)
      - [Create Replication Endpoints](administration/configuring_replication/create_replication_endpoints.md)
      - [Create Replication Rules](administration/configuring_replication/create_replication_rules.md)
      - [Manage Replications](administration/configuring_replication/manage_replications.md) 
  - [Vulnerability Scanning](administration/vulnerability_scanning/_index.md)
      - [Connect Harbor to Additional Vulnerability Scanners](administration/vulnerability_scanning/pluggable_scanners.md)
      - [Scan an Individual Image](administration/vulnerability_scanning/scan_individual_image.md)
      - [Scan All Images](administration/vulnerability_scanning/scan_all_images.md)
      - [Schedule Scans](administration/vulnerability_scanning/schedule_scans.md)
      - [Import Vulnerability Data to an Offline Harbor instance](administration/vulnerability_scanning/import_vulnerability_data.md)
      - [Configure System-Wide CVE Whitelists](administration/vulnerability_scanning/configure_system_whitelist.md)
  - [Garbage Collection](administration/garbage_collection/_index.md)
  - [Upgrading Harbor](administration/upgrade/_index.md)
      - [Upgrade Harbor and Migrate Data](administration/upgrade/upgrade_migrate_data.md)
      - [Roll Back an Upgrade](administration/upgrade/roll_back_upgrade.md)
  - [Manage the Harbor Instance](administration/manage_harbor/_index.md)
      - [Access Harbor Logs](administration/manage_harbor/access_logs.md)

## Working with Harbor Projects

This section describes how users with the developer, master, and project administrator roles manage and participate in Harbor projects.

- [Working with Harbor Projects](working_with_projects/_index.md)
  - [Project Creation](working_with_projects/project_overview.md)
      - [Create a Project](working_with_projects/create_projects.md)
      - [Assign Users to a Project](working_with_projects/add_users.md)
  - [Project Configuration](working_with_projects/project_configuration.md)
      - **[Access and Search Project Logs](working_with_projects/access_project_logs.md)**
      - [Create Robot Accounts](working_with_projects/create_robot_accounts.md)
      - [Configure Webhook Notifications](working_with_projects/configure_webhooks.md)
      - [Configure a Per-Project CVE Whitelist](working_with_projects/configure_project_whitelist.md)
      - **[Implementing Content Trust](working_with_projects/implementing_content_trust.md)**
  - [Working with Images, Tags, and Helm Charts](working_with_projects/working_with_images.md)
      - [Pulling and Pushing Images](working_with_projects/pulling_pushing_images.md)
      - **[Create Labels](working_with_projects/create_labels.md)**
      - **[Retag Images](working_with_projects/retagging_images.md) **
      - [Create Tag Retention Rules](working_with_projects/create_tag_retention_rules.md)
      - [Create Tag Immutability  Rules](working_with_projects/create_tag_immutability_rules.md)
      - [Manage Kubernetes Packages with Helm Charts](working_with_projects/managing_helm_charts.md)

## Build, Customize, and Contribute to Harbor

This section describes how developers can build from Harbor source code, customize their deployments, and contribute to the open-source Harbor project.

- [Build Harbor from Source Code](build_customize_contribute/compile_guide.md)
- [Developing the Harbor Frontend](build_customize_contribute/ui_contribution_get_started.md)
- [Customize the Harbor Look & Feel ](build_customize_contribute/customize_look_feel.md)
- [Developing for Internationalization](build_customize_contribute/developer_guide_i18n.md)
- [Using Make](build_customize_contribute/use_make.md)
- [View and test Harbor REST API via Swagger](build_customize_contribute/configure_swagger.md)
- [Registry Landscape](build_customize_contribute/registry_landscape.md)