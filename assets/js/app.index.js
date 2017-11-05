$(function () {
    console.log("Bootstraping app.index.js");
    $('#send').click(function () {
        console.log("Send clicked");
        let name = $('#name');
        let progress = $('.progress');
        let interval = setInterval(function () {
            console.log("Showing progress");
            progress.show()
        }, 100);
        name
            .attr("disabled", "disabled")
            .addClass('w3-disabled');
        $.ajax("/api/register", {
            cache: false,
            data: {
                Name: name.val()
            },
            type: "POST"
        }).done(function (msg) {
            window.location = "/dashboard";
        }).fail(function (xhr, status, errorThrown) {
            name.addClass('w3-text-red').removeAttr('disabled').removeClass('w3-disabled');
            $('#error').css('display', 'block').html(xhr.responseText);
            console.log(status);
        }).always(function (xhr, status) {
            clearInterval(interval);
            progress.hide();
        })
    })
});