# Copyright Project Harbor Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  This resource provides any keywords related to the Harbor private registry appliance

*** Variables ***
${create_project_button_xpath}  //clr-main-container//button[contains(., 'New Project')]
${project_name_xpath}  //*[@id='create_project_name']
${project_public_xpath}  //input[@name='public']/..//label
${project_save_css}  html body.no-scrolling harbor-app harbor-shell clr-main-container.main-container div.content-container div.content-area.content-area-override project div.row div.col-lg-12.col-md-12.col-sm-12.col-xs-12 div.row.flex-items-xs-between div.option-left create-project clr-modal div.modal div.modal-dialog div.modal-content div.modal-footer button.btn.btn-primary
${log_xpath}  //clr-main-container//clr-vertical-nav//a[contains(.,'Logs')]
${projects_xpath}  //clr-main-container//clr-vertical-nav//a[contains(.,'Projects')]
${project_replication_xpath}  //project-detail//a[contains(.,'Replication')]
${project_log_xpath}  //project-detail//li[contains(.,'Logs')]
${project_member_xpath}  //project-detail//li[contains(.,'Members')]

${create_project_CANCEL_button_xpath}  xpath=//button[contains(.,'CANCEL')]
${create_project_OK_button_xpath}  xpath=//button[contains(.,'OK')]
${delete_confirm_btn}  xpath=//confirmation-dialog//button[contains(.,'DELETE')]
${project_statistics_private_repository_icon}  xpath=//project/div/div/div[1]/div/statistics-panel/div/div[2]/div[1]/div[2]/div[2]/statistics/div/span[1]
${repo_delete_confirm_btn}  xpath=//clr-modal//button[2]
${repo_retag_confirm_dlg}  css=${modal-dialog}
${repo_delete_on_card_view_btn}  //clr-modal//button[contains(.,'DELETE')]
${delete_btn}  //button[contains(.,'Delete')]
${repo_delete_btn}  xpath=//repository//button[contains(.,'Delete')]
${project_delete_btn}  xpath=//list-project//clr-datagrid//button[contains(.,'Delete')]
${tag_delete_btn}  xpath=//tag-repository//clr-datagrid//button[contains(.,'Delete')]
${user_delete_btn}  xpath=/clr-dropdown-menu//button[contains(.,'Delete')]
