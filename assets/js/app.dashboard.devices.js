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
    $.ajax("/api/devices", {
        cache: false
    }).done(function (data) {
        const table = $("#mainTable");
        const tbody = table.find("tbody");
        for(let i = 0; i < data.length; i++) {
            let row = $("<tr>");
            row.append("<td>" + data[i].Mac + "</td>");
            row.append("<td>" + data[i].Ip + "</td>");
            row.append("<td>" + data[i].Name + "</td>");
            let records = $("<td>0</td>");
            row.append(records);
            records.load("/api/devices/" + data[i].Ip + "/records");
            let servers = $("<td>0</td>");
            row.append(servers);
            servers.load("/api/devices/" + data[i].Ip + "/servers");
            tbody.append(row)
        }
    })
}

$(() => {
    load_data();
});