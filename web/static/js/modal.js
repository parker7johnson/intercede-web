document.addEventListener('DOMContentLoaded', function() {
	const modal = document.getElementById('churchModal');
	const openBtn = document.getElementById('openChurchModalBtn');
	const closeBtn = document.getElementById('closeModalBtn');
	const input = document.getElementById('churchCodeInput');
	const saveBtn = document.getElementById('saveChurchCodeBtn');

	// Open modal
	if (openBtn) {
		openBtn.addEventListener('click', function() {
			modal.showModal();
		});
	}

	// Close modal function
	function closeModal() {
		modal.close();
		input.value = '';
		saveBtn.disabled = true;
	}

	// Close button
	if (closeBtn) {
		closeBtn.addEventListener('click', closeModal);
	}

	// Close on backdrop click
	modal.addEventListener('click', function(e) {
		const rect = modal.getBoundingClientRect();
		const isInDialog = (rect.top <= e.clientY && e.clientY <= rect.top + rect.height &&
			rect.left <= e.clientX && e.clientX <= rect.left + rect.width);
		if (!isInDialog) {
			closeModal();
		}
	});

	// Save button
	if (saveBtn) {
		saveBtn.addEventListener('click', function() {
			const code = input.value;
			console.log('Saving church code:', code);
			closeModal();
		});
	}

	// Enable/disable save button
	if (input) {
		input.addEventListener('input', function(e) {
			saveBtn.disabled = e.target.value.trim() === '';
		});
	}
});
