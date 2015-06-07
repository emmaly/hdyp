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
});
