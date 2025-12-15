*** Settings ***
Library    SeleniumLibrary
Test Teardown   Close All Browsers

*** Variables ***
${URL}    http://localhost/product/list
${BROWSER}    headlesschrome


*** Test Cases ***
ทดสอบ เข้าสู่ระบบไม่สำเร็จ ด้วย username ที่ไม่ถูกต้อง
    เข้าสู่เว็บไซต์ และตรวจสอบว่า redirect มาที่    /auth/login    login-page
    เข้าสู่ระบบไม่สำเร็จ    login-username-input    user_100    login-password-input    P@ssw0rd

*** Keywords ***
เข้าสู่เว็บไซต์ และตรวจสอบว่า redirect มาที่
    [Arguments]    ${target-url}    ${target-element-locator}
    Open Browser    url=${URL}    browser=${BROWSER}
    Wait Until Location Is Not    location=${URL}
    Location Should Contain    ${target-url}
    Page Should Contain Element    id:${target-element-locator}

เข้าสู่ระบบไม่สำเร็จ
    [Arguments]    ${username-input-locator}    ${username}    ${password-input-locator}    ${password}
    Wait Until Element Is Visible    id:${username-input-locator}
    Input Text    id:${username-input-locator}    ${username}
    Input Password    id:${password-input-locator}    ${password}
    Click Button    id:login-btn
    Alert Should Be Present    Invalid email or password.