var $ = require('jquery'),
    _ = require('underscore'),
    io = require('socket.io-client'),
    Handlebars = require('handlebars');

var KAIJU_URL="http://10.4.126.233:2714";

var Kaiju = function(options) {
    this.forum = options.forum;
    this.page = options.page;

    var socket = this.socket = io.connect(KAIJU_URL);

    this.socket.on('commentsFor', _.bind(this.handleCommentsFor, this));

    this._comments = { };

    this.commentTemplate = _.template(
        '<div class="comment" data-thread-id="<%= parentId %>">' +
        '<strong><%= user.name %></strong> (<%= user.email %>): <%= body %>' +
        '</div>")');

    this.threadTemplate = _.template('<div class="commentThread"></div>');
};

Kaiju.prototype.postComment = function(message) {
    console.log('postComment');
    this.socket.emit('postComment', JSON.stringify(_.extend({
        forum: this.forum,
        page: this.page
    }, message)));
};

Kaiju.prototype.getComments = function() {
    console.log('getComments', this.forum, this.page);
    this.socket.emit('getComments', JSON.stringify({
        forum: this.forum,
        page: this.page
    }));
};

Kaiju.prototype.handleCommentsFor = function(data) {
    var commentTemplate = this.commentTemplate,
        threadTemplate = this.threadTemplate;

    data = JSON.parse(data);

    if (_.isArray(data)) {

        var $commentSection = $('section.comments-section');

        $commentSection.empty();

        var $mainThread = $('<div class="comment-thread">');

        _.each(data, function(comment) {
            console.log(comment);
            $mainThread.append($(commentTemplate({
                user: {
                    name: comment.User.FullName,
                    email: comment.User.Email
                },
                body: comment.Body,
                parentId: comment.ParentID
            })));
        });

        $commentSection.append($mainThread);
    }
};

$(function() {
    var kaiju = new Kaiju({
        forum: "5346e494331583002c7de60e",
        page: "local_test_page"
    });

    kaiju.getComments();

    $('form.comment-form').on('submit', function(evt) {
        evt.preventDefault();
        evt.stopPropagation();

        var form = $('form.comment-form')[0];
        kaiju.postComment({
            user: form.user.value,
            email: form.email.value,
            body: form.body.value
        });
    });
});
