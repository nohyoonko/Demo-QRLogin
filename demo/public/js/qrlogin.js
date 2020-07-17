/* 임시로 random number 생성하는 함수 */
var generateRandom = function (min, max) {
    var ranNum = Math.floor(Math.random() * (max - min + 1)) + min;
    return ranNum;
}

/* QR Code 생성 -> text(랜덤값) 서버에서 생성해서 가지고 오기 */
var qrcode = new QRCode(document.getElementById("qrcode"), {
    width: 150,
    height: 150,
    colorDark: "#000000",
    colorLight: "#ffffff",
    correctLevel: QRCode.CorrectLevel.H
});

/* Modal Popup */
var modal = document.getElementById('myModal');
var btn = document.getElementById('myBtn');
var span = document.getElementsByClassName("close")[0];                                          

btn.onclick = function() {
    modal.style.display = "block";
    var randNum = generateRandom(0, 100);
    qrcode.makeCode(randNum);
    document.getElementById('hiddenNum').setAttribute('value', randNum);
}
span.onclick = function() {
    modal.style.display = "none";
    qrcode.clear();
    document.getElementById('hiddenNum').removeAttribute('value');
}
// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}

$("#qrcode > img").css({"margin": "5% auto"});