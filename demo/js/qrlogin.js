var box = document.getElementById("qrcode");
var qrcode = new QRCode(box, {
    width: 150,
    height: 150,
    colorDark: "#000000",
    colorLight: "#ffffff",
    correctLevel: QRCode.CorrectLevel.H
});

/* QR Code 생성 */
var btn = document.getElementById("myBtn");
btn.onclick = function(){
    $.post("/create", function (data) {
        alert("데이터 불러온 결과: " + data.randNum);
        qrcode.makeCode(data.randNum);
        var qrBox = document.getElementById("qr_box");
        qrBox.style.display = 'block';
        btn.style.display = 'none'; 
    });
}

/* QR Code 삭제 */
var cancel = document.getElementById("qrCancel");
cancel.onclick = function(){
    var qrBox = document.getElementById("qr_box");
    qrBox.style.display = 'none';
    btn.style.display = 'block';
    qrcode.clear();
}