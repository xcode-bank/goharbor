# -*- coding: utf-8 -*-

import time
import base
import v2_swagger_client
from v2_swagger_client.rest import ApiException

class Artifact(base.Base, object):
    def __init__(self):
        super(Artifact,self).__init__(api_type = "artifact")

    def list_artifacts(self, project_name, repo_name, **kwargs):
        client = self._get_client(**kwargs)
        return client.list_artifacts(project_name, repo_name)

    def get_reference_info(self, project_name, repo_name, reference, expect_status_code = 200, ignore_not_found = False,**kwargs):
        client = self._get_client(**kwargs)
        params = {}
        if "with_signature" in kwargs:
            params["with_signature"] = kwargs["with_signature"]
        if "with_tag" in kwargs:
            params["with_tag"] = kwargs["with_tag"]
        if "with_scan_overview" in kwargs:
            params["with_scan_overview"] = kwargs["with_scan_overview"]
        if "with_immutable_status" in kwargs:
            params["with_immutable_status"] = kwargs["with_immutable_status"]

        try:
            data, status_code, _ = client.get_artifact_with_http_info(project_name, repo_name, reference, **params)
            return data
        except ApiException as e:
            if e.status == 404 and ignore_not_found == True:
                return None
            else:
                raise Exception("Failed to get reference, {} {}".format(e.status, e.body))
        else:
            base._assert_status_code(expect_status_code, status_code)
            base._assert_status_code(200, status_code)
            return None

    def delete_artifact(self, project_name, repo_name, reference, expect_status_code = 200, expect_response_body = None, **kwargs):
        client = self._get_client(**kwargs)

        try:
             _, status_code, _ = client.delete_artifact_with_http_info(project_name, repo_name, reference)
        except ApiException as e:
            base._assert_status_code(expect_status_code, e.status)
            if expect_response_body is not None:
                base._assert_status_body(expect_response_body, e.body)
            return
        else:
            base._assert_status_code(expect_status_code, status_code)
            base._assert_status_code(200, status_code)

    def get_addition(self, project_name, repo_name, reference, addition, **kwargs):
        client = self._get_client(**kwargs)
        return client.get_addition_with_http_info(project_name, repo_name, reference, addition)

    def add_label_to_reference(self, project_name, repo_name, reference, label_id, **kwargs):
        client = self._get_client(**kwargs)
        label = v2_swagger_client.Label(id = label_id)
        return client.add_label_with_http_info(project_name, repo_name, reference, label)

    def copy_artifact(self, project_name, repo_name, _from, expect_status_code = 201, expect_response_body = None, **kwargs):
        client = self._get_client(**kwargs)

        try:
            data, status_code, _ = client.copy_artifact_with_http_info(project_name, repo_name, _from)
        except ApiException as e:
            base._assert_status_code(expect_status_code, e.status)
            if expect_response_body is not None:
                base._assert_status_body(expect_response_body, e.body)
            return
        else:
            base._assert_status_code(expect_status_code, status_code)
            base._assert_status_code(201, status_code)
            return data

    def create_tag(self, project_name, repo_name, reference, tag_name, expect_status_code = 201, ignore_conflict = False, **kwargs):
        client = self._get_client(**kwargs)
        tag = v2_swagger_client.Tag(name = tag_name)
        try:
            _, status_code, _ = client.create_tag_with_http_info(project_name, repo_name, reference, tag)
        except ApiException as e:
            if e.status == 409 and ignore_conflict == True:
                return
            else:
                raise Exception("Create tag error, {}.".format(e.body))
        else:
            base._assert_status_code(expect_status_code, status_code)

    def delete_tag(self, project_name, repo_name, reference, tag_name, expect_status_code = 200, **kwargs):
        client = self._get_client(**kwargs)
        try:
            _, status_code, _ = client.delete_tag_with_http_info(project_name, repo_name, reference, tag_name)
        except ApiException as e:
            base._assert_status_code(expect_status_code, e.status)
        else:
            base._assert_status_code(expect_status_code, status_code)

    def check_image_scan_result(self, project_name, repo_name, reference, expected_scan_status = "Success", **kwargs):
        timeout_count = 30
        scan_status=""
        while True:
            time.sleep(5)
            timeout_count = timeout_count - 1
            if (timeout_count == 0):
                break
            artifact = self.get_reference_info(project_name, repo_name, reference, **kwargs)
            scan_status = artifact.scan_overview['application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0'].scan_status
            if scan_status == expected_scan_status:
                return
        raise Exception("Scan image result is {}, not as expected {}.".format(scan_status, expected_scan_status))

    def check_reference_exist(self, project_name, repo_name, reference, ignore_not_found = False, **kwargs):
        artifact = self.get_reference_info( project_name, repo_name, reference, ignore_not_found=ignore_not_found, **kwargs)
        return {
            None: False,
        }.get(artifact, True)

    def waiting_for_reference_exist(self, project_name, repo_name, reference, ignore_not_found = True, period = 60, loop_count = 20, **kwargs):
        _loop_count = loop_count
        while True:
            print("Waiting for reference {} round...".format(_loop_count))
            _loop_count = _loop_count - 1
            if (_loop_count == 0):
                break
            artifact = self.get_reference_info(project_name, repo_name, reference, ignore_not_found=ignore_not_found, **kwargs)
            print("Returned artifact by get reference info:", artifact)
            if artifact  and artifact !=[]:
                return  artifact
            time.sleep(period)
        raise Exception("Reference is not exist {} {} {}.".format(project_name, repo_name, reference))
