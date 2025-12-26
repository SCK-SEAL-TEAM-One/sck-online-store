*** Settings ***
Library    SeleniumLibrary
Library    String
Library    Collections
Library    RPA.PDF
Library    OperatingSystem
Test Teardown   Cleanup Download Folder
Test Setup    Setup Folder For Download

*** Variables ***
${URL}    http://localhost/product/list
${BROWSER}    headlesschrome
${DOWNLOAD_DIR}    ${CURDIR}${/}temp_downloads

*** Test Cases ***
ทดสอบ สั่งซื้อสินค้า Balance Training Bicycle จัดส่งด้วย Kerry ชำระเงินด้วยบัตรเครดิต Visa สำเร็จ และตรวจสอบใบเสร็จ
    เข้าสู่เว็บไซต์ และตรวจสอบว่า redirect มาที่    /auth/login    login-page
    เข้าสู่ระบบ    login-username-input    user_1    login-password-input    P@ssw0rd
    เลือกดูสินค้า    product-card-name-1    Balance Training Bicycle
    ตรวจสอบรายละเอียดสินค้า    Balance Training Bicycle    SportsFun    4,314.60    43
    เพิ่มสินค้าลงตะกร้าและมีสินค้าจำนวน    1
    กลับไปที่หน้า homepage
    เลือกดูสินค้า    product-card-name-2    43 Piece dinner Set
    เพิ่มสินค้าลงตะกร้าและมีสินค้าจำนวน    2
    เพิ่มสินค้าลงตะกร้าและมีสินค้าจำนวน    2
    
    ${product1}=    Create Dictionary    id=1    brand=SportsFun    name=Balance Training Bicycle    points=43    price=4,314.60    qty=1
    ${product2}=    Create Dictionary    id=2    brand=CoolKidz     name=43 Piece dinner Set       points=9    price=931.62      qty=2
    ${product-list}=    Create List    ${product1}    ${product2}

    ตรวจสอบข้อมูลรายการสินค้าในตะกร้า และ Checkout    ${product-list}    5,246.22
    ใส่ที่อยู่จัดส่งสินค้า    
    ...    ณัฐพล    ศรีสมบัติ    
    ...    43/8 หมู่บ้านเปี่ยมสุข ถนนลาดพร้าว ซอย 63    กรุงเทพมหานคร
    ...    เขตวังทองหลาง    วังทองหลาง
    ...    10310    0891234567
    เลือกวิธีจัดส่งสินค้าเป็น    kerry
    ตรวจสอบค่าจัดส่งสินค้าของ Kerry เท่ากันกับ 50.00 บาท    kerry    50.00
    เลือกช่องทางการชำระเงินแบบ VISA Credit Card    Nattapon Srisombat    5123 4500 0000 0008    01/39    100
    ตรวจสอบราคารวมที่ต้องชำระเงิน ต้องเท่ากันกับ    5,296.22
    ยืนยัน OTP
    ตรวจสอบหมายเลขพัสดุว่าต้องขึ้นต้นด้วย    KR
    กดดาวน์โหลดไฟล์

    ${order-product1}=    Create Dictionary    id=1    brand=SportsFun    name=Balance Training Bicycle    points=43    price=4,314.60    qty=1
    ${order-product2}=    Create Dictionary    id=2    brand=CoolKidz     name=43 Piece dinner Set       points=4    price=465.81      qty=2
    ${order-product-list}=    Create List    ${order-product1}    ${order-product2}
    ตรวขสอบข้อมูลในไฟล์ PDF    sck    shuhari    KR    credit
    ...    ${order-product-list}
    ...    5,246.22    thai_post    5,296.22    52

*** Keywords ***
Setup Folder For Download
    Create Directory    ${DOWNLOAD_DIR} 
    Empty Directory     ${DOWNLOAD_DIR}

Cleanup Download Folder
    Remove Directory    ${DOWNLOAD_DIR}    recursive=True
    Close All Browsers

เข้าสู่เว็บไซต์ และตรวจสอบว่า redirect มาที่
    [Arguments]    ${target-url}    ${target-element-locator}
    ${prefs}=    Create Dictionary    
    ...    download.default_directory=${DOWNLOAD_DIR}
    ...    plugins.always_open_pdf_externally=${True}
    Open Browser    url=${URL}    browser=${BROWSER}    options=add_experimental_option("prefs", ${prefs})

    Delete All Cookies
    Execute Javascript    window.localStorage.clear();
    Execute Javascript    window.sessionStorage.clear();
    Location Should Contain    ${target-url}
    Page Should Contain Element    id:${target-element-locator}

เข้าสู่ระบบ
    [Arguments]    ${username-input-locator}    ${username}    ${password-input-locator}    ${password}
    Wait Until Element Is Visible    id:${username-input-locator}
    Input Text    id:${username-input-locator}    ${username}
    Input Password    id:${password-input-locator}    ${password}
    Click Button    id:login-btn
    Wait Until Location Is    ${URL}

เลือกดูสินค้า
    [Arguments]    ${card-name-locator}    ${expected-product-name}
    Wait Until Page Contains Element    id:${card-name-locator}
    Element Should Contain    id:${card-name-locator}    ${expected-product-name}
    Click Element    id:${card-name-locator}

ตรวจสอบรายละเอียดสินค้า
    [Arguments]    ${product-name}    ${product-brand}    ${product-thb-price}    ${product-point}
    Wait Until Element Is Visible    id:product-detail-product-name
    Element Text Should Be    id:product-detail-product-name    ${product-name}
    Element Text Should Be    id:product-detail-brand    ${product-brand}
    Element Text Should Be    id:product-detail-price-thb    ฿${product-thb-price}
    Element Text Should Be    id:product-detail-point    ${product-point} Points

เพิ่มสินค้าลงตะกร้าและมีสินค้าจำนวน
    [Arguments]    ${cart-qty}
    Click Button    id:product-detail-add-to-cart-btn
    Wait Until Element Contains    id:header-menu-cart-badge    ${cart-qty}

กลับไปที่หน้า homepage
    Click Element    id:header-logo-image
    Wait Until Location Is    ${URL}

ตรวจสอบข้อมูลรายการสินค้าในตะกร้า และ Checkout
    [Arguments]    ${product-list}    ${subtotal}
    Click Button    id:header-menu-cart-btn
    Wait Until Element Is Visible    id:product-1-price

    FOR    ${product}    IN    @{product-list}
        # Should Contain    ${text}    ${product.brand} - ${product.name}THB${product.price}${product.qty}
        Element Text Should Be    id:product-${product.id}-name    ${product.name}
        Element Text Should Be    id:product-${product.id}-price    ฿${product.price}
        Element Text Should Be    id:product-${product.id}-point    ${product.points} Points
        Element Attribute Value Should Be    id:product-${product.id}-quantity-input    value    ${product.qty}
    END

    Element Text Should Be    id:shopping-cart-subtotal-price    ฿${subtotal}
    Click Element    id:shopping-cart-checkout-btn

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

กดดาวน์โหลดไฟล์
    Click Button    id:download-order-summary-btn
    ${file_path}=    Wait For Download To Complete
    Set Test Variable    ${file_path}


ตรวขสอบข้อมูลในไฟล์ PDF
    [Arguments]    ${first_name}    ${last_name}    ${shipping_prefix}    ${payment_method_key}    
    ...    ${order_products}
    ...    ${subtotal}    ${shipping_method}    ${total_price}    ${receiving_points}

    ${pages}=    Get Text From PDF    ${file_path}
    ${pages}=    Evaluate    dict(${pages})
    ${text}=     Get From Dictionary    ${pages}    ${1}


    &{PAYMENT_METHOD}    Create Dictionary
    ...    credit=Credit Card / Debit Card
    ...    linepay=Line Pay

    &{DELIVERY_METHOD}    Create Dictionary
    ...    kerry=50.00
    ...    thai_post=50.00
    ...    lineman=id:100.00

    Should Contain    ${text}    Full Name: ${first_name} ${last_name}
    Should Contain    ${text}    Tracking Number: ${shipping_prefix}
    Should Contain    ${text}    Payment Method: ${PAYMENT_METHOD}[${payment_method_key}]

    FOR    ${product}    IN    @{order_products}
        Should Contain    ${text}    ${product.brand} - ${product.name}
        Should Contain    ${text}    ${product.price}
        Should Contain    ${text}    ${product.qty}
    END

    Should Contain    ${text}    Merchandise SubtotalTHB ${subtotal}
    Should Contain    ${text}    Shipping FeeTHB${DELIVERY_METHOD}[${shipping_method}]
    Should Contain    ${text}    Total PriceTHB ${total_price}
    Should Contain    ${text}    Receiving Point${receiving_points}


Wait For Download To Complete
    Wait Until Keyword Succeeds    20 sec    1 sec    Directory Should Not Be Empty    ${DOWNLOAD_DIR}
    @{files}=    List Files In Directory    ${DOWNLOAD_DIR}    *.pdf
    RETURN    ${DOWNLOAD_DIR}${/}${files}[0]
