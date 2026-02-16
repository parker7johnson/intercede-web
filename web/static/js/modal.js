const CHURCH_KEY = 'last_church';
document.addEventListener('DOMContentLoaded', function() {
  initializeButtonActions()
  updateChurchCode()

});

const updateChurchCode = () => {
  const churchCodeEl = document.getElementById('lastChurchCode');
  const lastChurch = localStorage.getItem(CHURCH_KEY);
  if (!lastChurch) {
    churchCodeEl.innerHTML = 'Enter Church Code here';
  } else if (churchCodeEl) {
    churchCodeEl.innerHTML = `Sharing with: ${lastChurch}`
  }


};

document.addEventListener('htmx:afterSwap', (event) => {
  initializeButtonActions();
  updateChurchCode()
});


const initializeButtonActions = () => {


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
}
