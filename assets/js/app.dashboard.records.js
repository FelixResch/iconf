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
    $.ajax("/api/records", {
        cache: false
    }).done(function (data) {
        const table = $("#mainTable");
        const tbody = table.find("tbody");
        for(let i = 0; i < data.length; i++) {
            const record = data[i];
            let row = $("<tr>");
            row.append("<td>" + record.Device +"</td>");
            row.append("<td>" + record.Type +"</td>");
            row.append("<td>" + record.Name +"</td>");
            row.append("<td>" + record.Description +"</td>");
            let link = $("<a><i class='fa fa-info-circle'></i></a>");
            link.attr("href", "/dashboard/records/" + record.Key);
            row.append($("<td></td>").append(link));
            tbody.append(row);
        }
    })
}

$(() => {
    load_data()
});