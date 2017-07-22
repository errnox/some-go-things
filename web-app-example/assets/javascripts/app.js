$(document).ready(function() {
  $('[data-toggle="tooltip"]').tooltip()

  // Hidables

  // Hide all hidables.
  $('.js-hidable .js-hidableBody').hide();
  // Style all hidables.
  $('.js-hidable .js-hidableBody').css({
    'position': 'absolute',
    'z-index': '9999',
    'padding': '10px',
    'border-radius': '4px',
    'background-color': '#FFFFFF',
    'border': '1px solid #CACACA',
    'box-shadow': '0 2px 2px rgba(0, 0, 0, 0.4)',
    'min-width': '200px',
  });
  // Add click logic.
  $('.js-hidable .js-hidableToggle').click(function() {
    var that = this;
    $(this).next('.js-hidableBody').toggle();
    $(this).next('.js-hidableBody').offset({
      // top: $(that).offset().top + $(that).height(),
      // left: $(that).offset().left - $(this).width() / 2,
      top: $(that).offset().top + $(that).height() * 2,
      left: $(that).offset().left,
    });
    if ($(this).parent('.js-hidable').hasClass('js-hidableAutoFocus')) {
      $(this).parent('.js-hidable').find('input').first().focus();
    }
    // Create (inisible) backdrop.
    var backdrop = $('<div class="js-hidableBackdrop">&nbsp;</div>').css({
      'width': '100%',
      'height': '100%',
      'position': 'fixed',
      'top': '0',
      'left': '0',
      'z-index': '9999',
      'background-color': 'rgba(0, 0, 0, 0.0)',
    });
    backdrop.click(function() {
      $(that).parent('.js-hidable').find('.js-hidableBody').hide();
      $(this).remove();
    });
    $(this).after(backdrop)
  });
  // Close buttons
  $('.js-hidableClose').click(function() {
    // Remove the backdrop
    $(this).parents('.js-hidable').first()
      .find('.js-hidableBackdrop').remove();
    // Hide the popup.
    $(this).parents('.js-hidableBody').hide();
  });
});
