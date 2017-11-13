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
    $('#activeDevices').load("/api/active/devices");
    $('#activeServers').load("/api/active/servers");
    $('#activeRecords').load("/api/active/records");
    $('#activeGames').load("/api/active/games");
}

$(() => {
    load_data();
});