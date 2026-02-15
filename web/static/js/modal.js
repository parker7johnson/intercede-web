const CHURCH_KEY = 'last_church';
document.addEventListener('DOMContentLoaded', function() {
  const modal = document.getElementById('churchModal');
  const openBtn = document.getElementById('openChurchModalBtn');
  const closeBtn = document.getElementById('closeModalBtn');
  const input = document.getElementById('churchCodeInput');
  const saveBtn = document.getElementById('saveChurchCodeBtn');

  updateChurchCode();

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
      localStorage.setItem(CHURCH_KEY, code);
      updateChurchCode();
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

const updateChurchCode = () => {
  const churchCodeEl = document.getElementById('lastChurchCode');
  const lastChurch = localStorage.getItem(CHURCH_KEY);
  if (!lastChurch) {
    churchCodeEl.innerHTML = 'Enter Church Code here';
  } else {
    churchCodeEl.innerHTML = `Sharing with: ${lastChurch}`
  }
  const prayerSubmitButton = document.querySelector('prayer-btn-submit');
  const praiseSubmitButton = document.querySelector('praise-btn-submit');
  if (praiseSubmitButton) {
    praiseSubmitButton.disabled = !(lastChurch === "" && code === null);
  }

  if (prayerSubmitButton) {
    prayerSubmitButton.disabled = lastChurch === "" || code === null;
  }

};
