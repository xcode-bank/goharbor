[Back to table of contents](../../_index.md)

----------

# Configure Database Authentication

In database authentication mode, user accounts are stored in the local database. By default, only the Harbor system administrator can create user accounts to add users to Harbor. You can optionally configure Harbor to allow self-registration.  

**IMPORTANT**: If you create users in the database, Harbor is locked in database mode. You cannot change to a different authentication mode after you have created local users.

1. Log in to the Harbor interface with an account that has Harbor system administrator privileges.
1. Under **Administration**, go to **Configuration** and select the **Authentication** tab.
1. Leave **Auth Mode** set to the default **Database** option.

   ![Database authentication](../../img/db_auth.png)
   
1. Optionally select the **Allow Self-Registration** check box.

   ![Enable self-registration](../../img/new_self_reg.png)
    
   If you enable self registration option, users can register themselves in Harbor. Self-registration is disabled by default. If you enable self-registration, unregistered users can sign up for a Harbor account by clicking **Sign up for an account** in the Harbor log in page.
    
    ![Enable self-registration](../../img/self-registration-login.png)

----------

[Back to table of contents](../../_index.md)
