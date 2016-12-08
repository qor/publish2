(function(factory) {
    if (typeof define === 'function' && define.amd) {
        // AMD. Register as anonymous module.
        define(['jquery'], factory);
    } else if (typeof exports === 'object') {
        // Node / CommonJS
        factory(require('jquery'));
    } else {
        // Browser globals.
        factory(jQuery);
    }
})(function($) {

    'use strict';
    var $document = $(document);
    var NAMESPACE = 'qor.publish2';
    var EVENT_ENABLE = 'enable.' + NAMESPACE;
    var EVENT_DISABLE = 'disable.' + NAMESPACE;
    var EVENT_CLICK = 'click.' + NAMESPACE;
    var EVENT_CHANGE = 'change.' + NAMESPACE;
    var EVENT_SELECTONE_SELECTED = 'qor.selectone.selected qor.selectone.unselected';
    var EVENT_REPLICATOR_ADDED = 'added.qor.replicator';

    // sharedable version input name, please change this if adjust name in template !!
    // <input name="QorResource.ColorVariations[0].SizeVariations[0].ShareableVersion" />
    var NAME_SHAREABLEVERSION = 'ShareableVersion';

    var CLASS_VERSION_LINK = '.qor-publish2__version';

    var VERSION_LIST = 'qor-table__inner-list';
    var VERSION_BLOCK = 'qor-table__inner-block';
    var CLASS_VERSION_LIST = '.' + VERSION_LIST;
    var CLASS_VERSION_BLOCK = '.' + VERSION_BLOCK;

    var CLASS_EVENT_ID = '.qor-pulish2__eventid';
    var CLASS_EVENT_INPUT = '.qor-pulish2__eventid-input';

    var CLASS_PUBLISH_READY = '[name="QorResource.PublishReady"]';
    var CLASS_SCHEDULED_STARTAT = '[name="QorResource.ScheduledStartAt"]';
    var CLASS_SCHEDULED_ENDAT = '[name="QorResource.ScheduledEndAt"]';
    var CLASS_PUBLISH_VERSIONNAME = '[name="QorResource.VersionName"]';
    var CLASS_PUBLISH_EVENTID = '[name="QorResource.ScheduledEventID"]';

    var CLASS_PUBLISH_ACTION = '.qor-pulish2__action';
    var CLASS_PUBLISH_ACTION_SHAREDVERSION = '.qor-pulish2__action-sharedversion';

    var CLASS_PUBLISH_ACTION_START = '.qor-pulish2__action-start';
    var CLASS_PUBLISH_ACTION_END = '.qor-pulish2__action-end';
    var CLASS_PUBLISH_ACTION_INPUT = '.qor-pulish2__action-input';
    var CLASS_PICKER_BUTTON = '.qor-action__picker-button';
    var CLASS_MEDIALIBRARY_TR = '>tbody>tr';

    var IS_MEDIALIBRARY = 'qor-table--medialibrary';
    var IS_SHOW_VERSION = 'is-showing';



    function QorPublish2(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, QorPublish2.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    QorPublish2.prototype = {
        constructor: QorPublish2,

        init: function() {
            this.actionType = this.options.element;
            this.bind();
            this.initActionTemplate();
        },

        bind: function() {
            $document
                .on(EVENT_CLICK, CLASS_VERSION_LINK, this.loadPublishVersion.bind(this))
                .on(EVENT_CHANGE, CLASS_PUBLISH_ACTION_INPUT, this.action.bind(this))
                .on(EVENT_SELECTONE_SELECTED, CLASS_EVENT_ID, this.eventidChanged.bind(this))
                .on(EVENT_REPLICATOR_ADDED, this.replicatorAdded.bind(this));
        },

        unbind: function() {
            $document
                .off(EVENT_CLICK, CLASS_VERSION_LINK, this.loadPublishVersion.bind(this))
                .off(EVENT_CHANGE, CLASS_PUBLISH_ACTION_INPUT, this.action.bind(this))
                .off(EVENT_SELECTONE_SELECTED, CLASS_EVENT_ID, this.eventidChanged.bind(this))
                .off(EVENT_REPLICATOR_ADDED, this.replicatorAdded.bind(this));
        },

        initActionTemplate: function() {
            if (!$(CLASS_PUBLISH_ACTION).closest('.qor-slideout').size()) {
                $(CLASS_PUBLISH_ACTION).prependTo($('.mdl-layout__content .qor-page__body')).show();
            }
            QorPublish2.initSharedVersion();
        },

        replicatorAdded: function(e, $element) {
            QorPublish2.generateSharedVersionLabel($element);
        },

        action: function(e) {
            var $element = $(e.target),
                isCheckbox = $element.prop('type') == 'checkbox',
                currentValue = $element.val(),
                $target = $(this.actionType[$element.data().actionType]),
                $checkboxLabel = $target.closest('label');

            if (!$target.size()) {
                return;
            }

            if (isCheckbox) {
                $target.prop('checked', $element.is(':checked'));
                $element.is(':checked') ? $checkboxLabel.addClass('is-checked') : $checkboxLabel.removeClass('is-checked');
            } else {
                $target.val(currentValue);
            }

        },

        eventidChanged: function(e, data) {
            if (data) {
                $(CLASS_EVENT_INPUT).val(data.primaryKey);
            } else {
                $(CLASS_EVENT_INPUT).val('');
            }

            this.updateDate(data, e.target);
            $(CLASS_EVENT_INPUT).trigger('change');
        },

        updateDate: function(data, element) {
            var $parent = $(element).closest(CLASS_PUBLISH_ACTION),
                $start = $parent.find(CLASS_PUBLISH_ACTION_START),
                $end = $parent.find(CLASS_PUBLISH_ACTION_END),
                $button = $parent.find(CLASS_PICKER_BUTTON),
                $input = $button.parent().find('input');

            if (data) {
                $start.val(data.ScheduledStartAt);
                $end.val(data.ScheduledEndAt);
                $button.hide();
                $input.attr('disabled', true);
            } else {
                $button.show();
                $input.attr('disabled', false).closest('.is-disabled').removeClass('is-disabled');
            }

            $start.trigger('change');
            $end.trigger('change');
        },

        loadPublishVersion: function(e) {
            var $target = $(e.target).parent("a"),
                url = $target.data().versionUrl,
                $table = $target.closest('table'),
                $tr = $target.closest('tr'),
                colspan = $tr.find('td').size(),
                isMediaLibrary = $table.hasClass(IS_MEDIALIBRARY),
                $list,
                $newRow = $('<tr class="' + VERSION_LIST + '"><td colspan="' + colspan + '"></td></tr>'),
                $version = $('<div class="' + VERSION_BLOCK + '"><div style="text-align: center;"><div class="mdl-spinner mdl-js-spinner is-active"></div></div></div>');

            if ($tr.hasClass(IS_SHOW_VERSION)) {
                $(CLASS_VERSION_LIST).remove();
                $table.find('tr').removeClass(IS_SHOW_VERSION);
                return false;
            }

            $(CLASS_VERSION_LIST).remove();
            $('table tr').removeClass(IS_SHOW_VERSION);

            $tr.addClass(IS_SHOW_VERSION);

            if (isMediaLibrary) {
                var $trs = $table.find(CLASS_MEDIALIBRARY_TR),
                    columnNum = parseInt($table.width() / 217),
                    currentNum = $trs.index($tr) + 1,
                    rows = Math.ceil($trs.size() / columnNum),
                    currentRow = Math.ceil(currentNum / columnNum);

                $tr = $($trs.get((columnNum * currentRow) - 1));
                if (!$tr.size()) {
                    $tr = $trs.last();
                }
                $newRow = $('<tr class="' + VERSION_LIST + '"><td></td></tr>');
                if (rows > 1) {
                    $newRow.width(217 * columnNum - 16);
                }
            }

            $tr.after($newRow);
            $list = $(CLASS_VERSION_LIST).find('td');

            $version.appendTo($list).trigger('enable');

            url && $.get(url, function(html) {
                $(CLASS_VERSION_BLOCK).html(html).trigger('enable');
            });

            return false;
        },

        destroy: function() {
            this.unbind();
            this.$element.removeData(NAMESPACE);
        }

    };

    QorPublish2.generateSharedVersionLabel = function($element) {
        var sharedVersion = $('[name="shared-version-checkbox"]').html(),
            $inputs = $('input[name$="' + NAME_SHAREABLEVERSION + '"]'),
            data = {},
            randomString;

        if ($element) {
            $inputs = $element.find('input[name$="' + NAME_SHAREABLEVERSION + '"]');
        }

        $inputs.each(function() {
            var $input = $(this),
                $field = $input.closest('.qor-fieldset'),
                $template;


            if ($field.hasClass('.qor-fieldset--new')) {
                return;
            }

            if ($element) {
                $field.find(CLASS_PUBLISH_ACTION_SHAREDVERSION).remove();
            }

            randomString = (Math.random() + 1).toString(36).substring(7);
            data.id = [NAME_SHAREABLEVERSION, randomString].join('_');
            $template = $(window.Mustache.render(sharedVersion, data));

            $template.find('input').on(EVENT_CLICK, function() {
                $(this).is(':checked') ? $input.val('true') : $input.val('');
            });

            if ($input.val() == "true") {
                $template.find('input').prop('checked', true);
            }

            $template.prependTo($field).trigger('enable');
            $input.closest('.qor-field').hide();
        });
    };

    QorPublish2.initSharedVersion = function() {
        if (!$(CLASS_PUBLISH_ACTION).size()) {
            return;
        }

        QorPublish2.generateSharedVersionLabel();

    };

    $.fn.qorSliderAfterShow.initSharedVersion = QorPublish2.initSharedVersion;

    $.fn.qorSliderAfterShow.initPublishForm = function() {
        var $action = $(CLASS_PUBLISH_ACTION),
            $types = $action.find('[data-action-type]'),
            element = QorPublish2.ELEMENT;

        if (!$action.size() || !$types.size()) {
            return;
        }

        $types.each(function() {
            var $this = $(this);
            $(element[$this.data().actionType]).closest('.qor-form-section').hide();
        });

        $(CLASS_PUBLISH_ACTION_INPUT).trigger(EVENT_CHANGE);

    };

    QorPublish2.DEFAULTS = {};

    QorPublish2.ELEMENT = {
        'scheduledstart': CLASS_SCHEDULED_STARTAT,
        'scheduledend': CLASS_SCHEDULED_ENDAT,
        'publishready': CLASS_PUBLISH_READY,
        'versionname': CLASS_PUBLISH_VERSIONNAME,
        'eventid': CLASS_PUBLISH_EVENTID
    };

    QorPublish2.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {

                if (/destroy/.test(options)) {
                    return;
                }

                $this.data(NAMESPACE, (data = new QorPublish2(this, options)));
            }

            if (typeof options === 'string' && $.isFunction(fn = data[options])) {
                fn.apply(data);
            }
        });
    };


    $(function() {
        var selector = '.qor-theme-publish2';

        var options = {};

        options['element'] = QorPublish2.ELEMENT;

        $(document).
        on(EVENT_DISABLE, function(e) {
            QorPublish2.plugin.call($(selector, e.target), 'destroy');
        }).
        on(EVENT_ENABLE, function(e) {
            QorPublish2.plugin.call($(selector, e.target), options);
        }).
        triggerHandler(EVENT_ENABLE);
    });

    return QorPublish2;

});