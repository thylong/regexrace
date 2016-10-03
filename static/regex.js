// DISCLAIMER: JS is not my thing, So the following can be really improved...

// Send the score at the end of the time to the /score handler.
function submitScore(score, username) {
    $.ajax({
        type: "POST",
        url: "/score",
        data: JSON.stringify({
            "best_score": score,
            "username": "user"+username
        }),
        dataType: "json",
        success: function(data) {
            console.log("Submitted score !");
        }
    });
}

// ticking starts the visual timer and trigger the callback when timing out.
function ticking(duration, display, callback) {
    var timer = duration, minutes, seconds;
    var interval = setInterval(function () {
        minutes = parseInt(timer / 60, 10);
        seconds = parseInt(timer % 60, 10);

        minutes = minutes < 10 ? "0" + minutes : minutes;
        seconds = seconds < 10 ? "0" + seconds : seconds;

        if (--timer < 0) {
            display.text("00:00");
            timer = duration;
            clearInterval(interval);
            submitScore(
                parseInt($("#score").text()),
                String(Math.floor((Math.random() * 100) + 1))
            );
            $("#timer-container").hide();
            $("#retry").show();
            $("input[name=regex]").attr("disabled", true);
            return;
        }
        display.text(minutes + ":" + seconds);

    }, 1000);
}

/********************************************************
******************* WHEN PAGE READY *********************
********************************************************/
jQuery(document).ready(function($) {
    localStorage.removeItem('trigger_timer');
    $("#retry").click(function() {
        location.reload();
    });

    display = $('#timer');
    duration = $("#timer").attr("data-duration");

    // Initiliaze timer with param in config.
    var timer = duration, minutes, seconds;
    minutes = parseInt(timer / 60, 10);
    seconds = parseInt(timer % 60, 10);
    minutes = minutes < 10 ? "0" + minutes : minutes;
    seconds = seconds < 10 ? "0" + seconds : seconds;
    display.text(minutes + ":" + seconds);

    $("input").keypress(function (e) {
        var trigger_timer = localStorage.getItem('trigger_timer');
        if( !trigger_timer ) {
            ticking(duration, display);
            localStorage.setItem('trigger_timer', true);
        }
        var key = e.which;
        if(key == 13) { // If the Key pressed is Enter, send call to /answer.
            var answer = {
                "qid": parseInt($("form").attr("id")),
                "regex": $("input[name=regex]").val(),
                "modifier": $("input[name=modifier]").val(),
            };

            $.ajax({
                type: "POST",
                url: "/answer",
                data: JSON.stringify(answer),
                dataType: "json",
                success: function(data) {
                    if (data.status == "success") {
                        console.log("success");
                        $("form>div").removeClass("has-error");
                        $("p#sentence").html(data.new_question.sentence);
                        $("input[name=regex]").val("");
                        $("input[name=modifier]").val("");
                        $("#score").text(parseInt($("#score").text()) + 1);
                        $("form").attr("id", data.new_question.qid);
                    } else {
                        console.log("fail");
                        $("form>div").addClass("has-error");
                    }
                    $("input[name=regex]").focus();
                }
            });
        }
    });
});
