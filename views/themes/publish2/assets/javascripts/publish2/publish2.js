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
  var CLASS_VERSION_LINK = '.qor-publish2__version';
  var VERSION_LIST = 'qor-table__inner-list';
  var VERSION_BLOCK = 'qor-table__inner-block';
  var CLASS_VERSION_LIST = '.' + VERSION_LIST;
  var CLASS_VERSION_BLOCK = '.' + VERSION_BLOCK;
  var CLASS_TABLE = '.qor-table';
  var IS_MEDIALIBRARY = 'qor-table--medialibrary';

  function QorPublish2(element, options) {
    this.$element = $(element);
    this.options = $.extend({}, QorPublish2.DEFAULTS, $.isPlainObject(options) && options);
    this.init();
  }

  QorPublish2.prototype = {
    constructor: QorPublish2,

    init: function () {
      this.bind();
    },

    bind: function () {
      $document
        .on(EVENT_CLICK, CLASS_VERSION_LINK, this.loadPublishVersion.bind(this));
    },

    unbind: function () {
      $document
        .off(EVENT_CLICK, CLASS_VERSION_LINK, this.loadPublishVersion.bind(this));
    },

    loadPublishVersion: function (e) {
      var $target = $(e.target),
          url = $target.data().versionUrl,
          $table = $target.closest('table'),
          $tr = $target.closest('tr'),
          hasVersion = $tr.next(CLASS_VERSION_LIST).size(),
          colspan = $tr.find('td').size(),
          isMediaLibrary = $table.hasClass(IS_MEDIALIBRARY),
          newRow = '<tr class="' + VERSION_LIST + '"><td colspan="' + colspan + '"></td></tr>',
          $version = $('<div class="' + VERSION_BLOCK + '"><div style="text-align: center;"><div class="mdl-spinner mdl-js-spinner is-active"></div></div></div>'),
          $list;

      if (hasVersion) {
        $(CLASS_VERSION_LIST).remove();
        return false;
      }

      $(CLASS_VERSION_LIST).remove();
      $tr.after(newRow);
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
