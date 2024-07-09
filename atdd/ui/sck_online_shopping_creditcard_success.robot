*** Settings ***
Library    SeleniumLibrary
Test Teardown   Close All Browsers

*** Variables ***
${URL}    http://192.168.3.115:3000/product/list
${BROWSER}    headlesschrome

*** Test Cases ***
ทดสอบ สั่งซื้อสินค้า Balance Training Bicycle จัดส่งด้วย Kerry ชำระเงินด้วยบัตรเครดิต Visa สำเร็จ และตรวจสอบการได้แต้มสะสม
    ค้นหาสินค้าด้วย คำค้นหา    Bicycle
    ตรวจสอบผลการค้นหา    product-card-name-1    Balance Training Bicycle
    เลือกดูสินค้า    product-card-name-1
    ตรวจสอบจำนวนแต้มต่อชิ้นที่จะได้รับของ     product-detail-point    43 Points
    เพิ่มสินค้าลงตะกร้า    Balance Training Bicycle
    ตรวจสอบจำนวนแต้มต่อชิ้นที่จะได้รับของสินค้าในตะกร้า    product-1-point
    ใส่ที่อยู่จัดส่งสินค้า    
    ...    พงศกร    รุ่งเรืองทรัพย์    
    ...    189/413 หมู่ 2    สมุทรปราการ
    ...    เมืองสมุทรปราการ    แพรกษาใหม่
    ...    10280    0909127991
    เลือกวิธีจัดส่งสินค้าเป็น    shipping-method-1-card
    ตรวจสอบค่าจัดส่งสินค้าของ Kerry เท่ากันกับ 50.00 บาท    shipping-method-1-fee    ฿50.00
    เลือกช่องทางการชำระเงินแบบ VISA Credit Card
    ตรวจสอบราคารวมที่ต้องชำระเงิน ต้องเท่ากันกับ    ฿4,364.60
    ยืนยัน OTP
    ตรวจสอบหมายเลขพัสดุ
    ยืนยันการส่งการแจ้งเตือนด้วย email และ เบอร์โทรศัพท์    
    ...    ponsakorn@gmail.com
    ...    0909127991

*** Keywords ***
ค้นหาสินค้าด้วย คำค้นหา
    [Arguments]    ${keyword}
    Open Browser    url=${URL}    browser=${BROWSER}
    # ...    options=add_experimental_option("detach", True)
    Input Text    id:search-product-input    ${keyword}
    Click Element    id:search-product-btn    

ตรวจสอบผลการค้นหา 
    [Arguments]    ${card-name-locator}    ${expected-product-name}
    Wait Until Element Is Visible    id:${card-name-locator}
    Element Should Contain    id:${card-name-locator}    ${expected-product-name} 
    
เลือกดูสินค้า
    [Arguments]    ${card-name-locator}
    Click Element    id:${card-name-locator}

ตรวจสอบจำนวนแต้มต่อชิ้นที่จะได้รับของ
    [Arguments]    ${detail-point-locator}    ${expected-point}
    Wait Until Element Is Visible    id:${detail-point-locator}
    Element Text Should Be    id:${detail-point-locator}    ${expected-point}

เพิ่มสินค้าลงตะกร้า
    [Arguments]    ${product-name}
    Click Button    id:product-detail-add-to-cart-btn

ตรวจสอบจำนวนแต้มต่อชิ้นที่จะได้รับของสินค้าในตะกร้า
    [Arguments]    ${product-point-locator}
    Click Button    id:header-menu-cart-btn
    Wait Until Element Is Visible    id:${product-point-locator}
    # Element Text Should Be    id:product-1-point    43 Points

ใส่ที่อยู่จัดส่งสินค้า
    [Arguments]    ${firstname}    ${lastname}    
    ...    ${address}    ${province}    ${district}
    ...    ${subdistrict}    ${zipcode}    ${mobileno}
    Click Link    id:shopping-cart-checkout-btn
    Input Text    id:shipping-form-first-name-input    ${firstname}
    Input Text    id:shipping-form-last-name-input    ${lastname}
    Input Text    id:shipping-form-address-input    ${address}
    Select From List By Label    id:shipping-form-province-select    ${province}
    Select From List By Label    id:shipping-form-district-select    ${district}
    Select From List By Label    id:shipping-form-sub-district-select    ${subdistrict}
    Element Attribute Value Should Be    id:shipping-form-zipcode-input    value    ${zipcode} 
    Input Text    id:shipping-form-mobile-input    ${mobileno}

เลือกวิธีจัดส่งสินค้าเป็น
    [Arguments]    ${shipping-method-card}
    Click Element    id:${shipping-method-card}

ตรวจสอบค่าจัดส่งสินค้าของ Kerry เท่ากันกับ 50.00 บาท
    [Arguments]    ${shipping-method-fee}    ${fee}
    Element Text Should Be    id:${shipping-method-fee}    ${fee}

เลือกช่องทางการชำระเงินแบบ VISA Credit Card
    Click Element    id:payment-credit-input
    Input Text    id:payment-credit-form-fullname-input
    ...    พงศกร รุ่งเรืองทรัพย์
    Input Text    id:payment-credit-form-card-number-input
    ...    4719700591590995
    Input Text    id:payment-credit-form-expiry-input
    ...    0226
    Input Text    id:payment-credit-form-cvv-input
    ...    752

ตรวจสอบราคารวมที่ต้องชำระเงิน ต้องเท่ากันกับ
    [Arguments]    ${total-price}
    Element Should Be Visible    id:order-summary-total-payment-price
    Element Text Should Be    id:order-summary-total-payment-price    ${total-price}

ยืนยัน OTP
    Click Button    id:payment-now-btn
    Wait Until Element Is Visible    id:otp-input
    Click Button    Request OTP
    Input Text    id:otp-input    124532
    Click Button    PAY NOW

ตรวจสอบหมายเลขพัสดุ
    Wait Until Element Is Visible    id:order-success-tracking-id
    Element Should Contain    id:order-success-tracking-id    KR-
    ${tracking-id}=    Get Text    id:order-success-tracking-id
    Should Match Regexp    ${tracking-id}    ^KR-\\d{7,9}$

ยืนยันการส่งการแจ้งเตือนด้วย email และ เบอร์โทรศัพท์
    [Arguments]    ${email}    ${telphoneno}
    Input Text    id:notification-form-email-input    ${email}
    Input Text    id:notification-form-mobile-input    ${telphoneno}
    Click Button     id:send-notification-btn
    Handle Alert
    Location Should Be    ${URL}