/* QR Code에 담을 random number 및 QR Code 생성 */
var rand = "";
for (var i = 0; i < 5; i++) {
    rand += (Math.floor((Math.random() * 10)));
}

var qrcode = new QRCode(document.getElementById("qrcode"), {
    text: rand,
    width: 128,
    height: 128,
    colorDark: "#000000",
    colorLight: "#ffffff",
    correctLevel: QRCode.CorrectLevel.H
});

$("#qrcode > img").css({"margin": "auto"});

/* Popup Layer */
$('.btn-example').click(function () {
    var $href = $(this).attr('href');
    layer_popup($href);
});

function layer_popup(el) {

    var $el = $(el); //레이어의 id를 $el 변수에 저장
    var isDim = $el.prev().hasClass('dimBg'); //dimmed 레이어를 감지하기 위한 boolean 변수

    isDim ? $('.dim-layer').fadeIn() : $el.fadeIn();

    var $elWidth = ~~($el.outerWidth()),
        $elHeight = ~~($el.outerHeight()),
        docWidth = $(document).width(),
        docHeight = $(document).height();

    // 화면의 중앙에 레이어를 띄운다.
    if ($elHeight < docHeight || $elWidth < docWidth) {
        $el.css({
            marginTop: -$elHeight / 2,
            marginLeft: -$elWidth / 2
        })
    } else {
        $el.css({
            top: 0,
            left: 0
        });
    }

    $el.find('a.btn-layerClose').click(function () {
        isDim ? $('.dim-layer').fadeOut() : $el.fadeOut(); // 닫기 버튼을 클릭하면 레이어가 닫힌다.
        return false;
    });

    $('.layer .dimBg').click(function () {
        $('.dim-layer').fadeOut();
        return false;
    });

}

var $qrnum = $('#qr_num')

$('#qr_form').on('submit', function(e){
    $.post('/qrcode', {
        qrnum: $qrnum.val()
    });
    $qrnum.val("");
    $qrnum.focus();
    return false;
});