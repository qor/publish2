(function (factory) {
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
})(function ($) {

  'use strict';
  var $document = $(document);
  var NAMESPACE = 'qor.publish2';
  var EVENT_ENABLE = 'enable.' + NAMESPACE;
  var EVENT_DISABLE = 'disable.' + NAMESPACE;
  var EVENT_CLICK = 'click.' + NAMESPACE;
  var EVENT_CHANGE = 'change.' + NAMESPACE;
  var EVENT_SELECTONE_SELECTED = 'qor.selectone.selected';

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
  var CLASS_PUBLISH_EVENTID = '[name="QorResource.ScheduleEventID"]';

  var CLASS_PUBLISH_ACTION_INPUT = '.qor-pulish2__action-input';
  var CLASS_MEDIALIBRARY_TR = '.qor-table--medialibrary>tbody>tr';

  var IS_MEDIALIBRARY = 'qor-table--medialibrary';
  var IS_SHOW_VERSION = 'is-showing';

  function QorPublish2(element, options) {
    this.$element = $(element);
    this.options = $.extend({}, QorPublish2.DEFAULTS, $.isPlainObject(options) && options);
    this.init();
  }

  QorPublish2.prototype = {
    constructor: QorPublish2,

    init: function () {
      this.actionType = {
        'scheduledstart' : CLASS_SCHEDULED_STARTAT,
        'scheduledend' : CLASS_SCHEDULED_ENDAT,
        'publishready': CLASS_PUBLISH_READY,
        'versionname': CLASS_PUBLISH_VERSIONNAME,
        'eventid': CLASS_PUBLISH_EVENTID
      };
      this.bind();
    },

    bind: function () {
      $document
        .on(EVENT_CLICK, CLASS_VERSION_LINK, this.loadPublishVersion.bind(this))
        .on(EVENT_CHANGE, CLASS_PUBLISH_ACTION_INPUT, this.action.bind(this))
        .on(EVENT_SELECTONE_SELECTED, CLASS_EVENT_ID, this.eventidChanged.bind(this));

    },

    unbind: function () {
      $document
        .off(EVENT_CLICK, CLASS_VERSION_LINK, this.loadPublishVersion.bind(this))
        .off(EVENT_CHANGE, CLASS_PUBLISH_ACTION_INPUT, this.action.bind(this));
    },

    action: function (e) {
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

    eventidChanged: function (e, data) {
      $(CLASS_EVENT_INPUT).val(data.primaryKey).trigger('change');
    },

    loadPublishVersion: function (e) {
      var $target = $(e.target),
          url = $target.data().versionUrl,
          $table = $target.closest('table'),
          $tr = $target.closest('tr'),
          colspan = $tr.find('td').size(),
          isMediaLibrary = $table.hasClass(IS_MEDIALIBRARY),
          $list,
          $newRow = $('<tr class="' + VERSION_LIST + '"><td colspan="' + colspan + '"></td></tr>'),
          $version = $('<div class="' + VERSION_BLOCK + '"><div style="text-align: center;"><div class="mdl-spinner mdl-js-spinner is-active"></div></div></div>');

      $(CLASS_VERSION_LIST).remove();
      $table.find('tr').removeClass(IS_SHOW_VERSION);

      $tr.addClass(IS_SHOW_VERSION);

      if (isMediaLibrary) {
        var $trs = $(CLASS_MEDIALIBRARY_TR),
            columnNum = parseInt($table.width() / 217),
            currentNum = $trs.index($tr) + 1,
            currentRow = Math.ceil(currentNum / columnNum);

        $tr = $($trs.get(( columnNum * currentRow ) - 1));
        if (!$tr.size()) {
          $tr = $trs.last();
        }
        $newRow = $('<tr class="' + VERSION_LIST + '" style="width: ' + (217 * columnNum - 16) + 'px"><td></td></tr>');
      }

      $tr.after($newRow);
      $list = $(CLASS_VERSION_LIST).find('td');

      $version.appendTo($list).trigger('enable');

      url && $.get(url, function (html) {
        $(CLASS_VERSION_BLOCK).html(html).trigger('enable');
      });

      return false;
    },

    destroy: function () {
      this.unbind();
      this.$element.removeData(NAMESPACE);
    }

  };

  $.fn.qorSliderAfterShow.initPublishForm = function () {

    if (!$('.qor-pulish2__action').size()) {
      return;
    }

    var classNames = [CLASS_SCHEDULED_STARTAT, CLASS_SCHEDULED_ENDAT, CLASS_PUBLISH_READY];

    for (var i = 0; i < classNames.length; i++) {
      $(classNames[i]).closest('.qor-form-section').hide();
    }

  };

  QorPublish2.DEFAULTS = {};

  QorPublish2.plugin = function (options) {
    return this.each(function () {
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


  $(function () {
    var selector = '.qor-theme-publish2';

    $(document).
      on(EVENT_DISABLE, function (e) {
        QorPublish2.plugin.call($(selector, e.target), 'destroy');
      }).
      on(EVENT_ENABLE, function (e) {
        QorPublish2.plugin.call($(selector, e.target));
      }).
      triggerHandler(EVENT_ENABLE);
  });

  return QorPublish2;

});
