*** Settings ***
Library    SeleniumLibrary
Test Teardown   Close All Browsers

*** Variables ***
${URL}    http://localhost/product/list
${BROWSER}    headlesschrome

*** Test Cases ***
ทดสอบ เคยข้าสู่ระบบและออกจากเว็บไซต์ แล้วกลับมาอีกครั้งเพื่อสั่งซื้อสินค้า Balance Training Bicycle จัดส่งด้วย Thai Post ชำระเงินด้วยบัตรเครดิต Visa สำเร็จ และตรวจสอบการได้แต้มสะสม
    เข้าสู่เว็บไซต์ และตรวจสอบว่า redirect มาที่    /auth/login    login-form-container
    เข้าสู่ระบบ    login-username-input    user_1    login-password-input    P@ssw0rd
    ออกจากเว็บไซต์ และเข้าสู่เว็บไซต์มาอีกครั้ง และตรวจสอบว่าไม่ redirect มาหน้า Login
    ค้นหาสินค้าด้วย คำค้นหา    Bicycle
    ตรวจสอบผลการค้นหา    product-card-name-1    Balance Training Bicycle
    เลือกดูสินค้า    product-card-name-1    Balance Training Bicycle
    ตรวจสอบรายละเอียดสินค้า    Balance Training Bicycle    SportsFun    4,314.60    43
    เพิ่มสินค้าลงตะกร้า
    ตรวจสอบข้อมูลสินค้าในตะกร้า และ Checkout    Balance Training Bicycle    SportsFun    4,314.60    43    4,314.60
    ใส่ที่อยู่จัดส่งสินค้า    
    ...    ณัฐพล    ศรีสมบัติ    
    ...    43/8 หมู่บ้านเปี่ยมสุข ถนนลาดพร้าว ซอย 63    กรุงเทพมหานคร
    ...    เขตวังทองหลาง    วังทองหลาง
    ...    10310    0891234567
    เลือกวิธีจัดส่งสินค้าเป็น    thai_post
    ตรวจสอบค่าจัดส่งสินค้าของ Kerry เท่ากันกับ 50.00 บาท    kerry    50.00
    เลือกช่องทางการชำระเงินแบบ VISA Credit Card    Nattapon Srisombat    5123 4500 0000 0008    01/39    100
    ตรวจสอบราคารวมที่ต้องชำระเงิน ต้องเท่ากันกับ    4,364.60
    ยืนยัน OTP
    ตรวจสอบหมายเลขพัสดุว่าต้องขึ้นต้นด้วย    TH

*** Keywords ***
เข้าสู่เว็บไซต์ และตรวจสอบว่า redirect มาที่
    [Arguments]    ${target-url}    ${target-element-locator}
    Open Browser    url=${URL}    browser=${BROWSER}
    Wait Until Location Is Not    location=${URL}
    Location Should Contain    ${target-url}
    Page Should Contain Element    id:${target-element-locator}

เข้าสู่ระบบ
    [Arguments]    ${username-input-locator}    ${username}    ${password-input-locator}    ${password}
    Wait Until Element Is Visible    id:${username-input-locator}
    Input Text    id:${username-input-locator}    ${username}
    Input Password    id:${password-input-locator}    ${password}
    Click Button    id:login-btn
    Wait Until Location Is    ${URL}


ออกจากเว็บไซต์ และเข้าสู่เว็บไซต์มาอีกครั้ง และตรวจสอบว่าไม่ redirect มาหน้า Login
    # Execute Javascript    window.open('${URL}');
    Go To    about:blank
    # "Re-enter" the website
    Go To    ${URL}
    Wait Until Location Is    ${URL}
    Wait Until Element Is Visible    id:search-product-input


ค้นหาสินค้าด้วย คำค้นหา
    [Arguments]    ${keyword}
    Input Text    id:search-product-input    ${keyword}
    Click Element    id:search-product-btn

ตรวจสอบผลการค้นหา
    [Arguments]    ${card-name-locator}    ${expected-product-name}
    Wait Until Element Is Visible    id:${card-name-locator}
    Element Should Contain    id:${card-name-locator}    ${expected-product-name} 

เลือกดูสินค้า
    [Arguments]    ${card-name-locator}    ${expected-product-name}
    Page Should Contain Element    id:${card-name-locator}
    Element Should Contain    id:${card-name-locator}    ${expected-product-name}
    Click Element    id:${card-name-locator}

ตรวจสอบรายละเอียดสินค้า
    [Arguments]    ${product-name}    ${product-brand}    ${product-thb-price}    ${product-point}
    Wait Until Element Is Visible    id:product-detail-product-name
    Element Text Should Be    id:product-detail-product-name    ${product-name}
    Element Text Should Be    id:product-detail-brand    ${product-brand}
    Element Text Should Be    id:product-detail-price-thb    ฿${product-thb-price}
    Element Text Should Be    id:product-detail-point    ${product-point} Points

เพิ่มสินค้าลงตะกร้า
    Click Button    id:product-detail-add-to-cart-btn
    Wait Until Element Contains    id:header-menu-cart-badge    text=1

ตรวจสอบข้อมูลสินค้าในตะกร้า และ Checkout
    [Arguments]    ${product-name}    ${product-brand}    ${product-thb-price}    ${product-point}    ${subtotal-price}
    Click Button    id:header-menu-cart-btn
    Wait Until Element Is Visible    id:product-1-price
    Element Text Should Be    id:product-1-name    ${product-name}
    Element Text Should Be    id:product-1-price    ฿${product-thb-price}
    Element Text Should Be    id:product-1-point    ${product-point} Points
    Element Text Should Be    id:shopping-cart-subtotal-price    ฿${subtotal-price}
    Click Element    id:shopping-cart-checkout-btn

ใส่ที่อยู่จัดส่งสินค้า
    [Arguments]    ${firstname}    ${lastname}    
    ...    ${address}    ${province}    ${district}
    ...    ${subdistrict}    ${zipcode}    ${phone-number}
    Input Text    id:shipping-form-first-name-input    ${firstname}
    Input Text    id:shipping-form-last-name-input    ${lastname}
    Input Text    id:shipping-form-address-input    ${address}
    Select From List By Label    id:shipping-form-province-select    ${province}
    Select From List By Label    id:shipping-form-district-select    ${district}
    Select From List By Label    id:shipping-form-sub-district-select    ${subdistrict}
    Element Attribute Value Should Be    id:shipping-form-zipcode-input    value    ${zipcode} 
    Input Text    id:shipping-form-mobile-input    ${phone-number}

เลือกวิธีจัดส่งสินค้าเป็น
    [Arguments]    ${method}
    &{DELIVERY_METHOD}    Create Dictionary
    ...    kerry=id:shipping-method-1-card
    ...    thai_post=id:shipping-method-2-card
    ...    lineman=id:shipping-method-3-card
    Click Element    ${DELIVERY_METHOD}[${method}]

ตรวจสอบค่าจัดส่งสินค้าของ Kerry เท่ากันกับ 50.00 บาท
    [Arguments]    ${method}    ${fee}
    &{DELIVERY_METHOD}    Create Dictionary
    ...    kerry=id:shipping-method-1-fee
    ...    thai_post=id:shipping-method-2-fee
    ...    lineman=id:shipping-method-3-fee
    Element Text Should Be    ${DELIVERY_METHOD}[${method}]    ฿${fee}

เลือกช่องทางการชำระเงินแบบ VISA Credit Card
    [Arguments]    ${credit-card-name}    ${credit-card-number}    ${credit-card-expired-date}       ${credit-card-cvv}
    Click Element    id:payment-credit-input
    Input Text    id:payment-credit-form-fullname-input    ${credit-card-name}
    Input Text    id:payment-credit-form-card-number-input    ${credit-card-number}
    Input Text    id:payment-credit-form-expiry-input    ${credit-card-expired-date}
    Input Text    id:payment-credit-form-cvv-input    ${credit-card-cvv}

ตรวจสอบราคารวมที่ต้องชำระเงิน ต้องเท่ากันกับ
    [Arguments]    ${total-price}
    Element Should Be Visible    id:order-summary-total-payment-price
    Element Text Should Be    id:order-summary-total-payment-price    ฿${total-price}

ยืนยัน OTP
    Click Button    id:payment-now-btn
    Wait Until Element Is Visible    id:otp-input
    Click Button    Request OTP
    Input Text    id:otp-input    124532
    Click Button    OK

ตรวจสอบหมายเลขพัสดุว่าต้องขึ้นต้นด้วย
    [Arguments]    ${shipping-prefix}
    Wait Until Element Is Visible    id:order-success-tracking-id
    Element Should Contain    id:order-success-tracking-id    ${shipping-prefix}-
    ${tracking-id}=    Get Text    id:order-success-tracking-id
    Should Match Regexp    ${tracking-id}    ^${shipping-prefix}-\\d{7,9}$

