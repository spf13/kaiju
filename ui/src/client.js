var $ = require('jquery'),
    _ = require('underscore'),
    io = require('socket.io-client'),
    Handlebars = require('handlebars');

var KAIJU_URL="http://10.4.126.233:2714";

var Kaiju = function(options) {
    this.forum = options.forum;
    this.page = options.page;

    var socket = this.socket = io.connect(KAIJU_URL);

    this.socket.on('commentsFor', this.handleCommentsFor);
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
    this.socket.emit('getComments', {
        forum: this.forum,
        page: this.page
    });
};

Kaiju.prototype.handleCommentsFor = function(data) {
    console.log("handle comments", data);
    if (_.isArray(data.comments)) {
        $commentSection = $('section.comments-section');

        $commentSection.empty();

        _.each(data.comments, function(comment) {
            console.log(comment);
            $commentSection.append('<div class="well"><strong>' +
                comment.user + "</strong> (" + comment.email + "): " + comment.body + "</div>");
        });
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
