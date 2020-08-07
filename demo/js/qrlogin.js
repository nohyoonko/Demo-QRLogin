/* QR Code 초기화 */
var qrDiv = document.getElementById("qrcode");
var qrcode = new QRCode(qrDiv, {
    width: 150,
    height: 150,
    colorDark: "#000000",
    colorLight: "#ffffff",
    correctLevel: QRCode.CorrectLevel.H
});
var verify = ""; //인증 여부를 확인하는 변수

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
    var jsonData = {
        "randNum": qrNum
    }
    $.ajax({
        type: "DELETE",
        url: "/create",
        data: JSON.stringify(jsonData),
        contentType: 'application/json; charset=utf-8'
    });

    var qrBox = document.getElementById("qrBox");
    qrBox.style.display = 'none';
    createBtn.style.display = 'block';
    qrcode.clear();
}

/* Mobile 역할을 하는 prompt */
qrDiv.onclick = function () {
    var cred = prompt(qrNum + "의 credential을 입력하세요."); //취소하면 null 반환
    var mobileData = {
        "randNum": qrNum,
        "credential": cred
    }
    verify = "";

    if (cred != null) {
        $.ajax({
            type: "POST",
            url: "/mobile",
            data: JSON.stringify(mobileData), //jsonObj -> jsonStr
            contentType: 'application/json; charset=utf-8',
            success: function (response) {
                verify = response.verify;
                if (verify == "success") {
                    alert("인증에 성공하였습니다!");
                }
            },
            error: function (response, status, error) {
                if (response.status == "401") { //http.StatusUnauthorized
                    alert("인증에 실패하였습니다!");
                    var json = JSON.parse(response.responseText); //jsonStr -> jsonObj
                    verify = json.verify;
                } else {
                    alert(error);
                }
            }
        });
    }
}

/* QR Code 확인(success or fail) */
var checkBtn = document.getElementById("qrCheck");
checkBtn.onclick = function () {
    //인증 과정을 거친 후에만 동작하도록
    if (verify != "") {
        var jsonData = {
            "randNum": qrNum
        }
        $.ajax({
            type: "POST",
            url: "/check",
            data: JSON.stringify(jsonData),
            contentType: 'application/json; charset=utf-8',
            success: function (response) {
                if (response.QRstatus == "success") {
                    window.location.href = "/success"
                } else if (response.QRstatus == "fail") {
                    window.location.href = "/"
                }
            },
            error: function (error) {
                alert("Error!" + error);
            }
        });
    }
}