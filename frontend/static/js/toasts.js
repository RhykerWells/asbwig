function sendSuccessToast(message) {
    const container = document.getElementById("toastContainer");

    const toast = Object.assign(document.createElement("div"), {
        className: `toast align-items-center bg-success`,
        role: "alert",
    });
    toast.setAttribute("aria-live", "assertive");
    toast.setAttribute("aria-atomic", "true");
    toast.innerHTML = `
    <div class="toast-header bg-success">
        <strong class="me-auto">Success</strong>
        <button type="button" class="btn-close btn-close-white" data-bs-dismiss="toast" aria-label="Close"></button>
    </div>
    <div class="toast-body">
        ${message}
    </div>
    `;
    container.appendChild(toast)

    const bsToast = new bootstrap.Toast(toast, { delay: 5000 });
    bsToast.show();

    toast.addEventListener("hidden.bs.toast", () => toast.remove());
}

function sendFailureToast(message) {
    const container = document.getElementById("toastContainer");

    const toast = Object.assign(document.createElement("div"), {
        className: `toast align-items-center bg-danger`,
        role: "alert",
    });
    toast.setAttribute("aria-live", "assertive");
    toast.setAttribute("aria-atomic", "true");
    toast.innerHTML = `
    <div class="toast-header bg-danger">
        <strong class="me-auto">Failure</strong>
        <button type="button" class="btn-close btn-close-white" data-bs-dismiss="toast" aria-label="Close"></button>
    </div>
    <div class="toast-body">
        ${message} (ask support if you need help)
    </div>
    `;
    container.appendChild(toast)

    const bsToast = new bootstrap.Toast(toast, { delay: 3000 });
    bsToast.show();

    toast.addEventListener("hidden.bs.toast", () => toast.remove());
}