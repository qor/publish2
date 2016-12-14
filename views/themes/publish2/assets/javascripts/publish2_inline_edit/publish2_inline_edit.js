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

    var location = window.location;
    var $document = $(document);
    var NAMESPACE = 'qor.filter';
    var EVENT_ENABLE = 'enable.' + NAMESPACE;
    var EVENT_DISABLE = 'disable.' + NAMESPACE;
    var EVENT_CLICK = 'click.' + NAMESPACE;
    var CLASS_SCHEDULED_TIME = '.qor-filter__scheduled-time';
    var CLASS_SEARCH_PARAM = '[data-search-param]';
    var CLASS_FILTER_SELECTOR = '.qor-filter__dropdown';
    var CLASS_FILTER_TOGGLE = '.qor-filter-toggle';

    function QorFilterTime(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, QorFilterTime.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    QorFilterTime.prototype = {
        constructor: QorFilterTime,

        init: function() {
            this.bind();
            var $element = this.$element;

            this.$scheduleTime = $element.find(CLASS_SCHEDULED_TIME);
            this.$searchButton = $element.find(this.options.button);
            this.$trigger = $element.find(this.options.trigger);

            this.publishReadyOn = $('#qor-publishready__on').data().label;
            this.publishReadyOff = $('#qor-publishready__off').data().label;

            this.initActionTemplate();

        },

        bind: function() {
            var options = this.options;

            this.$element.
            on(EVENT_CLICK, options.trigger, this.show.bind(this)).
            on(EVENT_CLICK, options.clear, this.clear.bind(this)).
            on(EVENT_CLICK, options.button, this.search.bind(this));

            $document.on(EVENT_CLICK, this.close);
        },

        unbind: function() {
            var options = this.options;
            this.$element.
            off(EVENT_CLICK, options.trigger, this.show.bind(this)).
            off(EVENT_CLICK, options.clear, this.clear.bind(this)).
            off(EVENT_CLICK, options.button, this.search.bind(this));
        },

        initActionTemplate: function() {
            var scheduleTime = this.getUrlParameter('publish_scheduled_time'),
                publishReady = this.getUrlParameter('publish_ready'),
                $trigger = this.$trigger,
                $selectorLabel = $trigger.find('.qor-selector-label'),
                $publishReadyLabel = $trigger.find('.qor-publishready-label'),
                $publishreadyInput = $('#qor-publishready__on');


            if (publishReady == '') {
                $('[name="QorResource.PublishReady"]').prop('checked', false);
            } else {
                $publishReadyLabel.parent().show();
                if (publishReady === 'true') {
                    $publishreadyInput.prop('checked', true);
                    $publishReadyLabel.html(this.publishReadyOn);
                } else {
                    $publishreadyInput.prop('checked', false);
                    $publishReadyLabel.html(this.publishReadyOff);
                }
                $publishReadyLabel.before('<i class="material-icons qor-selector-clear" data-type="publishready">clear</i>');
            }

            if (scheduleTime != '') {
                this.$scheduleTime.val(scheduleTime);
                $selectorLabel.parent().show();
                $selectorLabel.html(scheduleTime);
                $selectorLabel.before('<i class="material-icons qor-selector-clear">clear</i>');
            }
        },

        show: function() {
            this.$element.find(CLASS_FILTER_SELECTOR).toggle();
        },

        close: function(e) {
            var $target = $(e.target),
                $filter = $(CLASS_FILTER_SELECTOR),
                filterVisible = $filter.is(':visible'),
                isInFilter = $target.closest(CLASS_FILTER_SELECTOR).size(),
                isInToggle = $target.closest(CLASS_FILTER_TOGGLE).size(),
                isInModal = $target.closest('.qor-modal').size();

            if (filterVisible && (isInFilter || isInToggle || isInModal)) {
                return;
            }
            $filter.hide();
        },

        clear: function(e) {
            var $element = $(e.target),
                $trigger = this.$trigger,
                $label = $trigger.find('.qor-selector-label'),
                $publishReadyLabel = $trigger.find('.qor-publishready-label');

            if ($element.data().type) {
                $publishReadyLabel.html('').parent().hide();
                $('[name="QorResource.PublishReady"]').prop('checked', false);
            } else {
                $label.parent().hide();
                this.$scheduleTime.val('');
            }

            $element.remove();
            this.$searchButton.click();
            return false;

        },

        getUrlParameter: function(name) {
            var search = location.search;
            name = name.replace(/[\[]/, '\\[').replace(/[\]]/, '\\]');
            var regex = new RegExp('[\\?&]' + name + '=([^&#]*)');
            var results = regex.exec(search);
            return results === null ? '' : decodeURIComponent(results[1].replace(/\+/g, ' '));
        },

        updateQueryStringParameter: function(key, value, uri) {
            var href = uri || location.href,
                escapedkey = String(key).replace(/[\\^$*+?.()|[\]{}]/g, '\\$&'),
                re = new RegExp('([?&])' + escapedkey + '=.*?(&|$)', 'i'),
                separator = href.indexOf('?') !== -1 ? '&' : '?';

            if (href.match(re)) {
                return href.replace(re, '$1' + key + '=' + value + '$2');
            } else {
                return href + separator + key + '=' + value;
            }
        },

        search: function() {
            var $searchParam = this.$element.find(CLASS_SEARCH_PARAM),
                uri,
                _this = this;

            if (!$searchParam.size()) {
                return;
            }

            $searchParam.each(function() {
                var $this = $(this),
                    searchParam = $this.data().searchParam,
                    hasCheckedLabel = $this.find('[name="QorResource.PublishReady"]').filter(':checked').size(),
                    val = $this.val();
                if (searchParam == 'publish_ready') {
                    if (hasCheckedLabel) {
                        val = !!Number($this.find('input:checked').val());
                    } else {
                        val = '';
                    }
                }

                uri = _this.updateQueryStringParameter(searchParam, val, uri);
            });
            window.location.href = uri;
        },

        destroy: function() {
            this.unbind();
            this.$element.removeData(NAMESPACE);
        }
    };

    QorFilterTime.DEFAULTS = {
        trigger: false,
        button: false,
        clear: false
    };

    QorFilterTime.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }

                $this.data(NAMESPACE, (data = new QorFilterTime(this, options)));
            }

            if (typeof options === 'string' && $.isFunction(fn = data[options])) {
                fn.apply(data);
            }
        });
    };

    $(function() {
        var selector = '[data-toggle="qor.filter.time"]';
        var options = {
            trigger: 'a.qor-filter-toggle',
            button: '.qor-filter__button-search',
            clear: '.qor-selector-clear'
        };

        $(document).
        on(EVENT_DISABLE, function(e) {
            QorFilterTime.plugin.call($(selector, e.target), 'destroy');
        }).
        on(EVENT_ENABLE, function(e) {
            QorFilterTime.plugin.call($(selector, e.target), options);
        }).
        triggerHandler(EVENT_ENABLE);
    });

    return QorFilterTime;

});