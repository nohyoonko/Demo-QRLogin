/* QR Code 초기화 */
var qrDiv = document.getElementById("qrcode");
var qrcode = new QRCode(qrDiv, {
    width: 150,
    height: 150,
    colorDark: "#000000",
    colorLight: "#ffffff",
    correctLevel: QRCode.CorrectLevel.H
});

/* QR Code 생성 */
var qrNum;
var createBtn = document.getElementById("myBtn");
createBtn.onclick = function () {
    $.post("/create", function (data) {
        qrNum = data.randNum;
        qrcode.makeCode(qrNum);
        qrDiv.setAttribute("title", "")
        var qrBox = document.getElementById("qrBox");
        qrBox.style.display = 'block';
        createBtn.style.display = 'none';
    });
}

/* QR Code 삭제 */
var cancelBtn = document.getElementById("qrCancel");
cancelBtn.onclick = function () {
    var qrBox = document.getElementById("qrBox");
    qrBox.style.display = 'none';
    createBtn.style.display = 'block';
    qrcode.clear();
    //취소할 때 생성한 random number 삭제 필요
}

/* Mobile 역할을 하는 prompt */
qrDiv.onclick = function () {
    var cred = prompt(qrNum+"의 credential을 입력하세요.");
    var mobileData = {
        "randNum": qrNum,
        "credential": cred
    }
    $.ajax({
        type: "POST",
        url: "/mobile",
        data: JSON.stringify(mobileData),
        contentType: 'application/json; charset=utf-8',
        success: function () {
            alert("Success!");
        },
        error: function () {
            alert("Error!");
        }
    });
}

/* QR Code 확인(success or fail) */
var checkBtn = document.getElementById("qrCheck");
checkBtn.onclick = function () {
    var jsonData = {
        "randNum": qrNum
    }
    $.ajax({
        type: "POST",
        url: "/check",
        data: JSON.stringify(jsonData),
        contentType: 'application/json; charset=utf-8',
        success: function (response) {
            if(response.QRstatus == "success"){
                window.location.href="/success"
            }
            else { //fail
                window.location.href="/"
            }
        },
        error: function (error) {
            alert("Error!"+error);
        }
    });
}