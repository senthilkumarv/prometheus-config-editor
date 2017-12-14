let editor = CodeMirror.fromTextArea(document.getElementById("code"), {
    lineNumbers: true,
    mode: "yaml",
    gutters: ["CodeMirror-lint-markers"],
    lint: true
});

$.get("/load")
    .done(function (content) {
        editor.setValue(content)
    });

$("#save-config").click(() => {
    $("#save-config").attr("disabled", "true");
    $("#apply-config-change").show();
    $.post("/save", editor.getValue())
        .done((content) => {
            $("#resultModalLabel").text(content["Status"]);
            $("#modalContent").text(content["Details"]);
            $("#resultModal").modal('show')
        })
        .fail((content) => {
            let result = content.responseJSON;
            $("#resultModalLabel").text(result["Status"]);
            $("#modalContent").text(result["Details"] + ". " + result["Error"]);
            $("#apply-config-change").hide();
            $("#resultModal").modal('show')
        })
        .always(() => {
            $("#save-config").removeAttr("disabled");
        });
});

$("#apply-config-change").click(() => {
    $.post("/apply")
        .done(() => {
            $("#resultModal").modal('hide');
        })
        .fail((content) => {
            $("#resultModalLabel").text("Failed to apply");
            let responseJSON = content.responseJSON;
            $("#modalContent").text(responseJSON["Details"] + ". " +  responseJSON["Error"]);
        });
});

