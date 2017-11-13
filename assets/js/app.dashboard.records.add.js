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

const state = Object();

function init() {
    state.tab = 1;
    state.data = Object();
    state.ui = Object();
    const name = $('#name');
    const type = $('#type');
    const tabs = [$('#tab-1'), $('#tab-2'), $('#tab-3')];
    const tabContents = [$('#tab-content-1'), $('#tab-content-2')];
    const back = $('#back');
    const next = $('#next');
    const finish = $('#finish');
    next.click(function () {
        if(state.tab === 1) {
            if(name.val().length > 0) {
                state.data.name = name.val();
                tabs[0].removeClass('w3-border-red').addClass('w3-border-green');
                tabContents[0].hide();
                tabContents[1].show();
                tabs[1].removeClass("w3-hover-text-gray").removeClass('w3-text-gray').removeClass('w3-border-gray').addClass('w3-border-red');
                back.show();
                state.tab = 2;
            } else {
                name.addClass("w3-text-red")
            }
        } else {
            if(type.val().length > 0) {
                state.data.type = type.val();
                tabs[1].removeClass('w3-border-red').addClass('w3-border-green');
                tabContents[1].hide();
                if(state.data.type === 'A') {
                    tabContents[2] = $('#tab-content-3-a');
                    state.ui.host = $('#host-a');
                    state.data.usePort = false;
                } else if(state.data.type === 'CNAME') {
                    tabContents[2] = $('#tab-content-3-cname');
                    state.ui.host = $('#host-cname');
                    state.data.usePort = false;
                } else if(state.data.type === 'SRV') {
                    tabContents[2] = $('#tab-content-3-srv');
                    state.ui.host = $('#host-srv');
                    state.ui.port = $('#port-srv');
                    state.data.usePort = true;
                }
                tabContents[2].show();
                tabs[2].removeClass("w3-hover-text-gray").removeClass('w3-text-gray').removeClass('w3-border-gray').addClass('w3-border-red');
                finish.show();
                next.hide();
                state.tab = 3;
            } else {
                type.addClass("w3-text-red")
            }
        }
    });
    back.click(function () {
       if(state.tab === 2) {
           tabs[1].removeClass('w3-border-red').addClass('w3-border-green');
           tabContents[1].hide();
           tabContents[0].show();
           tabs[0].removeClass('w3-border-green').addClass('w3-border-red');
           back.hide();
           state.tab = 1;
       } else {
           tabs[2].removeClass('w3-border-red').addClass('w3-border-green');
           tabContents[2].hide();
           tabContents[1].show();
           tabs[1].removeClass('w3-border-green').addClass('w3-border-red');
           next.show();
           finish.hide();
           state.tab = 2;
       }
    });
    finish.click(function () {
        if(state.tab === 3) {
            if (state.ui.host.val().length > 0 && (!state.data.usePort || state.ui.port.val().length > 0)) {
                state.data.host = state.ui.host.val();
                if(state.data.usePort) {
                    state.data.port = state.ui.port.val()
                }
                $.ajax("/api/records", {
                    method: "POST",
                    data: JSON.stringify(state.data)
                }).done(function (data) {
                    if(data.State) {
                        document.location = "/dashboard/records/" + data.Identifier
                    } else {
                        alert(data)
                    }
                }).fail(function (jqwhr, error, other) {
                    console.log(jqwhr, error, other)
                })
            } else {
                alert("Please enter all data")
            }
        }
    })
}

$(() => init());