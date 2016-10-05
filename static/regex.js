// DISCLAIMER: JS is not my thing, So the following can be really improved...

// Send the score at the end of the time to the /score handler.
function submitScore() {
    score = parseInt($("#score").text());
    username = $("#username").val();

    if (username !== "" && username !== null) {
        $.ajax({
            type: "POST",
            url: "/score",
            data: JSON.stringify({
                "best_score": score,
                "username": username,
                "token": localStorage.getItem('regexrace_token'),
            }),
            dataType: "json",
            beforeSend : function(xhr) {
                xhr.setRequestHeader("Authorization", "Bearer "+localStorage.getItem('regexrace_token'));
            },
            success: function(data) {
                $("form.score>div").removeClass("has-error");
                $('#scoreModal').modal('toggle');
                console.log("Submitted score !");
            },
            error: function(data) {
                $("form.score>div").addClass("has-error");
                console.log("Submitted score !");
            }
        });
    } else {
        $("form.score>div").addClass("has-error");
    }
}

// Call /auth to get a token.
function getToken() {
    $.ajax({
        type: "GET",
        url: "/auth",
        success: function(data) {
            localStorage.setItem('regexrace_token', data.token);
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
            $("#submitModalButton").show();
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
    localStorage.removeItem('regexrace_timer');
    localStorage.removeItem('regexrace_token');
    $("#retry").click(function() {
        location.reload();
    });
    $("#submitScore").click(function() {
        submitScore();
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
        var regexrace_timer = localStorage.getItem('regexrace_timer');
        if( !regexrace_timer ) {
            ticking(duration, display);
            getToken();
            localStorage.setItem('regexrace_timer', true);
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
                beforeSend : function(xhr) {
                    xhr.setRequestHeader("Authorization", "Bearer "+localStorage.getItem('regexrace_token'));
                },
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
