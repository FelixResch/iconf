let sidebar = $('#mySidebar');
let overlay = $('#myOverlay');

function menu_open() {
    if(sidebar.css('display') === 'block') {
        $(sidebar, overlay).css('display', 'none')
    } else {
        $(sidebar, overlay).css('display', 'block')
    }
}

function menu_close() {
    $(sidebar, overlay).css('display', 'none')
}

function load_data() {
    const key = $('#key').html();
    $.ajax('/api/records/' + key, {cache: false})
        .done(function (data) {
            let short = data.Type + ':' + data.Name + ' &rarr; ' + data.Description;
            $('#recordKey').html(short);
            $('#detailRecordKey').text(data.Key);
            $('#detailRecordType').text(data.Type);
            $('#detailRecordName').text(data.Name);
            $('#detailRecordDescription').text(data.Description);
            $('#modalRecordDescription').html(short);
            let modal = $('#deleteModal');
            $('#closeModal').click(function () {
                modal.css('display', 'none')
            });
            $('#modalBack').click(function () {
                modal.css('display', 'none')
            });
            $('#modalDelete').click(function () {
                //TODO show progress bar after 100ms
                $.ajax("/api/records/" + key, {
                    method: "DELETE"
                }).done(function (data) {
                    if(data.State) {
                        document.location = "/dashboard/records"
                    } else {
                        $('#error').text(
                            'Could not delete ' +
                            data.Identifier +
                            ': ' + data.Reason
                        )
                    }
                }).fail(function (err) {
                    console.log(err)
                });
            });
            $('#delete').click(function () {
                modal.css('display', 'block')
            });
        })
}

$(() => {
    load_data()
});