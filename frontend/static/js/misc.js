function disableScreenPassthrough() {
    const overlay = document.createElement("div");
	overlay.className = "page-overlay d-flex";
    overlay.innerHTML = `<div class="spinner-border text-success" role="status"></div>`

    Object.assign(overlay.style, {
        position: "fixed",
        inset: "0",
        top: "0",
        left: "0",
        width: "100%",
        height: "100%",
        background: "rgba(0,0,0,0.3)",
        zIndex: "9999",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        cursor: "none"
    });
	document.body.prepend(overlay);
}

function reloadPage(delaySeconds) {
    const delay = (delaySeconds || 0) * 1000;

    setTimeout(() => {
        location.reload();
    }, delay);
}