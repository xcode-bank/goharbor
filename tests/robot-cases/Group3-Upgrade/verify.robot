*** Settings ***
Documentation  Harbor BATs
Resource  ../../resources/Util.robot
Default Tags  Nightly

*** Variables ***
${HARBOR_URL}  https://${ip}
${HARBOR_ADMIN}  admin

*** Test Cases ***
Test Case - Upgrade Verify
    [Tags]  1.8-latest
    ${data}=  Load Json From File  ${CURDIR}${/}data.json
    Run Keyword  Verify User  ${data}
    Run Keyword  Verify Project  ${data}
    Run Keyword  Verify Member Exist  ${data}
    Run Keyword  Verify Robot Account Exist  ${data}
    Run Keyword  Verify User System Admin Role  ${data}
    Run Keyword  Verify Endpoint  ${data}
    Run Keyword  Verify Replicationrule  ${data}
    Run Keyword  Verify System Setting  ${data}
    Run Keyword  Verify Image Tag  ${data}

Test Case - Upgrade Verify
    [Tags]  1.9-latest
    ${data}=  Load Json From File  ${CURDIR}${/}data.json
    Run Keyword  Verify User  ${data}
    Run Keyword  Verify Project  ${data}
    Run Keyword  Verify Member Exist  ${data}
    Run Keyword  Verify Robot Account Exist  ${data}
    Run Keyword  Verify Project-level Allowlist  ${data}
    Run Keyword  Verify Webhook  ${data}
    Run Keyword  Verify Tag Retention Rule  ${data}
    Run Keyword  Verify User System Admin Role  ${data}
    Run Keyword  Verify Endpoint  ${data}
    Run Keyword  Verify Replicationrule  ${data}
    Run Keyword  Verify Interrogation Services  ${data}
    Run Keyword  Verify System Setting  ${data}
    Run Keyword  Verify System Setting Allowlist  ${data}
    Run Keyword  Verify Image Tag  ${data}
    Run Keyword  Verify Trivy Is Default Scanner

Test Case - Upgrade Verify
    [Tags]  1.10-latest
    ${data}=  Load Json From File  ${CURDIR}${/}data.json
    Run Keyword  Verify User  ${data}
    Run Keyword  Verify Project  ${data}
    Run Keyword  Verify Member Exist  ${data}
    Run Keyword  Verify Robot Account Exist  ${data}
    Run Keyword  Verify Project-level Allowlist  ${data}
    Run Keyword  Verify Webhook  ${data}
    Run Keyword  Verify Tag Retention Rule  ${data}
    Run Keyword  Verify Tag Immutability Rule  ${data}
    Run Keyword  Verify User System Admin Role  ${data}
    Run Keyword  Verify Endpoint  ${data}
    Run Keyword  Verify Replicationrule  ${data}
    Run Keyword  Verify Interrogation Services  ${data}
    Run Keyword  Verify System Setting  ${data}
    Run Keyword  Verify System Setting Allowlist  ${data}
    Run Keyword  Verify Image Tag  ${data}
    Run Keyword  Verify Clair Is Default Scanner

Test Case - Upgrade Verify
    [Tags]  2.0-latest
    ${data}=  Load Json From File  ${CURDIR}${/}data.json
    Run Keyword  Verify User  ${data}
    Run Keyword  Verify Project  ${data}
    Run Keyword  Verify Member Exist  ${data}
    Run Keyword  Verify Robot Account Exist  ${data}
    Run Keyword  Verify Project-level Allowlist  ${data}
    Run Keyword  Verify Webhook For 2.0  ${data}
    Run Keyword  Verify Tag Retention Rule  ${data}
    Run Keyword  Verify Tag Immutability Rule  ${data}
    Run Keyword  Verify User System Admin Role  ${data}
    Run Keyword  Verify Endpoint  ${data}
    Run Keyword  Verify Replicationrule  ${data}
    Run Keyword  Verify Interrogation Services  ${data}
    Run Keyword  Verify System Setting  ${data}
    Run Keyword  Verify System Setting Allowlist  ${data}
    Run Keyword  Verify Image Tag  ${data}
    Run Keyword  Verify Trivy Is Default Scanner
    Run Keyword  Verify Artifact Index  ${data}

Test Case - Upgrade Verify
    [Tags]  2.1-latest
    ${data}=  Load Json From File  ${CURDIR}${/}data.json
    Run Keyword  Verify User  ${data}
    Run Keyword  Verify Project  ${data}  check_content_trust=${false}
    Run Keyword  Verify Member Exist  ${data}
    Run Keyword  Verify Robot Account Exist  ${data}
    Run Keyword  Verify Project-level Allowlist  ${data}
    Run Keyword  Verify Webhook For 2.0  ${data}
    Run Keyword  Verify Tag Retention Rule  ${data}
    Run Keyword  Verify Tag Immutability Rule  ${data}
    Run Keyword  Verify User System Admin Role  ${data}
    Run Keyword  Verify Endpoint  ${data}
    Run Keyword  Verify Replicationrule  ${data}
    Run Keyword  Verify Interrogation Services  ${data}
    Run Keyword  Verify System Setting  ${data}
    Run Keyword  Verify System Setting Allowlist  ${data}
    Run Keyword  Verify Image Tag  ${data}
    Run Keyword  Verify Trivy Is Default Scanner
    Run Keyword  Verify Artifact Index  ${data}
    Run Keyword  Verify Proxy Cache Image Existence  ${data}
    Run Keyword  Verify Distributions  ${data}
    Run Keyword  Verify P2P Preheat Policy  ${data}
