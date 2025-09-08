document.addEventListener('DOMContentLoaded', () => {
    saveRoleOptionsMulti()
});


function saveRoleOptionsMulti() {
    document.querySelectorAll('form[id$="Form"]').forEach(form => {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            const data = form.querySelector('input[type="hidden"]');
            const ModAction = String(data.id).slice(0, -5);
            const Roles = JSON.parse(data.value);
            const body = {ModAction, Roles};

            try {
                const response = await fetch(form.getAttribute('action'), {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(body),
                });
                const result = await response.json();

                if (result.success) {
                    sendSuccessToast(result.message);
                } else {
                    sendFailureToast(result.error);
                }
            } catch (err) {
                console.error(err);
                sendFailureToast('An unexpected error occured. Please try again\nerr: ${err}');
            }
        })
    });
};
