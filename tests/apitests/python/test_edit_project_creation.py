from __future__ import absolute_import
import unittest

from testutils import CLIENT
from testutils import TEARDOWN
from library.project import Project
from library.user import User
from library.configurations import Configurations

class TestProjects(unittest.TestCase):
    @classmethod
    def setUp(self):
        conf = Configurations()
        self.conf= conf

        project = Project()
        self.project= project

        user = User()
        self.user= user

    @classmethod
    def tearDown(self):
        print "Case completed"

    @unittest.skipIf(TEARDOWN == False, "Test data should be remain in the harbor.")
    def test_ClearData(self):
        print "Clear trace"
        #1. Delete project(PA);
        self.project.delete_project(TestProjects.project_edit_project_creation_id, **TestProjects.USER_edit_project_creation_CLIENT)

        #2. Delete user(UA);
        self.user.delete_user(TestProjects.user_edit_project_creation_id, **TestProjects.ADMIN_CLIENT)

    def testEditProjectCreation(self):
        """
        Test case:
            Edit Project Creation
        Test step & Expectation:
            1. Create a new user(UA);
            2. Set project creation to "admin only";
            3. Create a new project(PA) by user(UA), and fail to create a new project;
            4. Set project creation to "everyone";
            5. Create a new project(PA) by user(UA), success to create a project.
        Tear down:
            1. Delete project(PA);
            2. Delete user(UA);
        """
        admin_user = "admin"
        admin_pwd = "Harbor12345"
        url = CLIENT["endpoint"]
        user_edit_project_creation_password = "Aa123456"
        TestProjects.ADMIN_CLIENT=dict(endpoint = url, username = admin_user, password =  admin_pwd)

        #1. Create a new user(UA);
        TestProjects.user_edit_project_creation_id, user_edit_project_creation_name = self.user.create_user_success(user_password = user_edit_project_creation_password, **TestProjects.ADMIN_CLIENT)

        TestProjects.USER_edit_project_creation_CLIENT=dict(endpoint = url, username = user_edit_project_creation_name, password = user_edit_project_creation_password)

        #2. Set project creation to "admin only";
        self.conf.set_configurations_of_project_creation_restriction_success("adminonly", **TestProjects.ADMIN_CLIENT)

        #3. Create a new project(PA) by user(UA), and fail to create a new project;
        self.project.create_project(metadata = {"public": "false"}, expect_status_code = 403,
            expect_response_body = "Only system admin can create project", **TestProjects.USER_edit_project_creation_CLIENT)

        #4. Set project creation to "everyone";
        self.conf.set_configurations_of_project_creation_restriction_success("everyone", **TestProjects.ADMIN_CLIENT)

        #5. Create a new project(PA) by user(UA), success to create a project.
        TestProjects.project_edit_project_creation_id, _ = self.project.create_project(metadata = {"public": "false"}, **TestProjects.USER_edit_project_creation_CLIENT)


if __name__ == '__main__':
    unittest.main()