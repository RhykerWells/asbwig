document.addEventListener('DOMContentLoaded', () => {
    updateRoleOptionsMulti()
});

function updateRoleOptionsMulti() {
    document.querySelectorAll('.dropDownRoleCheckbox').forEach(checkbox => {
        checkbox.addEventListener('change', () => {
            const container = checkbox.closest('.input-group'); // Drop down container. Makes sure we only select objects for the appropriate role select.
            const checked = Array.from(container.querySelectorAll('.dropDownRoleCheckbox:checked'));
            const ids = checked.map(c => c.value);
            const names = checked.map(c => c.nextSibling.textContent.trim());
            const hiddenInput = container.querySelector('input[type=hidden]');


            hiddenInput.value = JSON.stringify(ids);

            const label = container.querySelector('span[id$="Label"]');
            let displayText = "Select roles";
            const joined = names.join(', ');
            if (names.length > 0) {
                displayText = joined;
                if (joined.length > 30) {
                    displayText = `${checked.length} Selected`;
                }
            }
            label.textContent = displayText
        });
    });
}