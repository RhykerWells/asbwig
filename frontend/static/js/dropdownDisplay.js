document.addEventListener('DOMContentLoaded', () => {
	updateRoleOptionsSingle()
	updateRoleOptionsMulti()
	updateChannelOptionsSingle()
});

/**
* Updates single-role dropdowns.
* Adds click event listeners to each element with the class 'dropDownRoleSingleItem'.
* When a role is clicked, the hidden input's value is updated and the label text is adjusted.
* Disabled items are ignored.
*/
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

/**
* Updates multi-role checkboxes.
* Adds change event listeners to each checkbox with the class 'dropDownRoleCheckbox'.
* Updates the hidden input with selected role IDs and updates the label text.
* If multiple roles are selected, the label either lists them or shows the count if the text is long.
*/
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

/**
* Updates single-channel dropdowns.
* Adds click event listeners to each element with the class 'channelListItem'.
* When a channel is clicked, updates the hidden input and updates the label text.
* Disabled items are ignored.
*/
function updateChannelOptionsSingle() {
	document.querySelectorAll('.channelListItem').forEach(item => {
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
			
			hiddenInput.value = value;
			
			const label = container.querySelector('span[id$="Label"]');
			let displayText = "Select channel";
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