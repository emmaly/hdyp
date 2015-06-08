$(function(){
  $(document).keyup(function(e){
    if (e.keyCode == 191) {
      $("input#word").focus().select();
    }
  });
  $("input#word").keyup(function(e){
    if (e.keyCode == 13) {
      document.location.href = "/" + $(this).val();
    }
    e.stopPropagation();
  });
  $("input#word:not([noautofocus])").focus().select();

  convertSVGtoInline();
});

function convertSVGtoInline() {
  // Replace all SVG images with inline SVG
  // Source: http://stackoverflow.com/questions/24933430/img-src-svg-changing-the-fill-color
  // ... but has been altered here since then.
  jQuery('img.svg').each(function(){
      var $img = jQuery(this);
      var imgID = $img.attr('id');
      var imgClass = $img.attr('class');
      var imgURL = $img.attr('src');

      jQuery.get(imgURL, function(data) {
          // Get the SVG tag, ignore the rest
          var $svg = jQuery(data).find('svg');
          // Add replaced image's ID to the new SVG
          if(typeof imgID !== 'undefined') {
              $svg = $svg.attr('id', imgID);
          }
          // Add replaced image's classes to the new SVG
          if(typeof imgClass !== 'undefined') {
              $svg = $svg.attr('class', imgClass+' replaced-svg');
          }
          // Remove any invalid XML tags as per http://validator.w3.org
          $svg = $svg.removeAttr('xmlns:a');
          // Set preserveAspectRatio
          $svg = $svg.attr('preserveAspectRatio', 'xMidYMid');
          // Set no viewbox
          $svg = $svg.attr('viewbox', 'none');
          // Replace image with new SVG
          $img.replaceWith($svg);
      }, 'xml');
  });
}
