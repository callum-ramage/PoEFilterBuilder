var updateEditor = null;

require(['plugins'], function () {
    require(['jquery', 'PoE/Helpers', 'itemfilter'], function ($, PoEHelpers, itemFilter) {
        $(document).ready(function () {
            var textarea = $('.filter-text');
            var text = textarea.val();
            var showPlaceholder = text == '' && !$('.filter-text-wrapper').hasClass('show-error');
            itemFilter.init('editor', false, {"blocks":["Show","Hide","Continue"],"conditions":["AlternateQuality","AnyEnchantment","AreaLevel","BaseType","BlightedMap","Class","Corrupted","CorruptedMods","DropLevel","ElderItem","ElderMap","EnchantmentPassiveNode","EnchantmentPassiveNum","FracturedItem","GemLevel","GemQualityType","HasEnchantment","HasExplicitMod","HasInfluence","Height","Identified","ItemLevel","LinkedSockets","MapTier","Mirrored","Prophecy","Quality","Rarity","Replica","ShapedMap","ShaperItem","SocketGroup","Sockets","StackSize","SynthesisedItem","Width"],"actions":["CustomAlertSound","DisableDropSound","EnableDropSound","MinimapIcon","PlayAlertSound","PlayAlertSoundPositional","PlayEffect","SetBackgroundColor","SetBorderColor","SetFontSize","SetTextColor"],"rarities":["Normal","Magic","Rare","Unique"],"gemQualityTypes":["Superior","Divergent","Anomalous","Phantasmal"],"colors":["Red","Green","Blue","Brown","White","Yellow","Cyan","Grey","Orange","Pink","Purple"],"shapes":["Circle","Diamond","Hexagon","Square","Star","Triangle","Cross","Moon","Raindrop","Kite","Pentagon","UpsideDownHouse"],"sockets":["R","G","B","A","D","W"],"influences":["Shaper","Elder","Crusader","Hunter","Redeemer","Warlord","None"],"operationTypes":["=","!","<=",">=","<",">","=="]});
            var editor = itemFilter.getEditor();
            editor.setValue(showPlaceholder ? textarea.attr('placeholder') : text, -1);
            $('.loading-filter').hide();
            $('#editor').show();

            updateEditor = function(content) {
                editor.setValue(content);
                // textarea.val(editor.getSession().getValue());
            }

            // Delete placeholder text if it displays
            editor.on('focus', function () {
                if (showPlaceholder) {
                    editor.setValue('', -1);
                    showPlaceholder = false;
                }
            });

            // Update text area value if editor changes
            editor.on('change', function () {
                if (!showPlaceholder)
                    textarea.val(editor.getSession().getValue());
            });

            $('.filter-text-element .error li').each(function (i, v) {
                if ($(v).find('.navigate-error').length == 0) {
                    var lineNumberRegex = /(^\d{1,})|^\w+\s(\d{1,})|^\W+\s(\d{1,})/g;
                    var lineNumbers = lineNumberRegex.exec($(v).text());
                    var number = 0;
                    for (var index in lineNumbers) {
                        if (!lineNumbers[index])
                            continue;
                        number = parseInt(lineNumbers[index], 10);
                        if (number && number > 0)
                            break;
                    }
                    if (number && number > 0)
                        $(v).append(`&nbsp;<a class="navigate-error" data-line="` + number + `">` + PoEHelpers.translate('Go to Line') + `</a>`);
                }
            });


            $('.filter-text-element .navigate-error, .item-filter-backend-error .navigate-error').click(function () {
                var linenumber = $(this).data('line');
                if (!linenumber) return;
                linenumber = parseInt(linenumber);
                // Navigate error line
                if (linenumber > 0) {
                    editor.gotoLine(linenumber, 0, true);
                }
            });
        });
    });
});