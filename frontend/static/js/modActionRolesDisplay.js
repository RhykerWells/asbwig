document.addEventListener('DOMContentLoaded', () => {
	updateRoleOptionsSingle()
	updateRoleOptionsMulti()
});

function updateRoleOptionsSingle() {
	document.querySelectorAll('.dropDownRoleSingleItem').forEach(item => {
		item.addEventListener('click', () => {
			if (item.classList.contains('disabled')) {
				e.stopPropagation();
				e.preventDefault();
				return;
			}

			const container = item.closest('.input-group'); // Drop down container. Makes sure we only select objects for the appropriate role select.
			const name = item.textContent.trim();
			const value = item.getAttribute('data-value');
			const hiddenInput = container.querySelector('input[type=hidden]');
			
			hiddenInput.value = JSON.stringify(value);
			
			const label = container.querySelector('span[id$="Label"]');
			let displayText = "Select role";
			if (value) {
				if (name.length > 30) {
					displayText = "1 Selected";
				} else {
					displayText = name;
				}
			}
			label.textContent = displayText
		});
	});
}

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