var $ = require('jquery'),
    _ = require('underscore'),
    io = require('socket.io-client'),
    Handlebars = require('handlebars');

var Kaiju = function(options) {
    this.forum = options.forum;
    this.page = options.page;
    this.$el = $(options.selector);

    var socket = this.socket = io.connect(options.url);

    this.socket.on('commentsFor', _.bind(this.handleCommentsFor, this));

    this._comments = { };

    this.commentTemplate = _.template(
        '<div class="comment" data-id="<%= id %>">' +
        '<div class="comment-header"><strong><%= user.name %></strong> (<%= user.email %>):</div>' +
        '<div class="comment-body"><%= body %></div>' +
        '<div class="comment-actions"><small><a class="add-comment" data-in-reply-to="<%= id %>" style="cursor:pointer;">Reply to this</a></small></div>' +
        '</div>")');

    this.threadTemplate = _.template('<div class="comment-thread"></div>');

    this.commentForm = this.$el.siblings('form.comment-form');
    this.addCommentLink = this.$el.siblings('a.add-comment');

    this.commentForm.on('submit', _.bind(this.onSubmitCommentForm, this));
    this.addCommentLink.on('click', _.bind(this.onClickShowCommentForm, this));
};

Kaiju.prototype.connect = function() {
    this.getComments();
};

Kaiju.prototype.postComment = function(message) {
    this.socket.emit('postComment', JSON.stringify(_.extend({
        forum: this.forum,
        page: this.page
    }, message)));
};

Kaiju.prototype.getComments = function() {
    this.socket.emit('getComments', JSON.stringify({
        forum: this.forum,
        page: this.page
    }));
};

Kaiju.prototype.handleCommentsFor = function(data) {
    var commentTemplate = this.commentTemplate,
        threadTemplate = this.threadTemplate,
        self = this;

    data = JSON.parse(data);

    if (_.isArray(data)) {

        _.each(data, function(comment) {
            self._comments[comment.Id] = comment;
        });

        var $mainThread = $(threadTemplate());
        this.$el.empty().append($mainThread);

        _.each(data, function(comment) {
            self.renderComment(comment, $mainThread);
        });

    }
};

Kaiju.prototype.renderComment = function(comment, $thread) {
    var commentTemplate = this.commentTemplate,
        parent = comment.Parent,
        $comment = $(commentTemplate({
            user: {
                name: comment.User.FullName,
                email: comment.User.Email
            },
            body: comment.Body,
            id: comment.Id,
            parentId: comment.Parent || "null"
        }));

    if (parent) {
        var $parentComment = this.$el.find("div.comment[data-id='" + parent + "']"),
            $thread = $parentComment.find('div.comment-thread');

        console.log("parent", parent, $parentComment, $parentComment.length);
    }

    $thread.append($comment);
    $comment.find('a.add-comment').on('click', _.bind(this.onClickShowCommentForm, this));
};

Kaiju.prototype.onSubmitCommentForm = function(evt) {
    evt.stopPropagation();
    evt.preventDefault();

    var form = this.commentForm[0];

    console.log({
        user: form.user.value,
        email: form.email.value,
        body: form.body.value,
        parent: form.parent.value
    });

    this.postComment({
        user: form.user.value,
        email: form.email.value,
        body: form.body.value,
        parent: form.parent.value
    });
};

Kaiju.prototype.onClickShowCommentForm = function(evt) {
    evt.stopPropagation();
    evt.preventDefault();

    this.showCommentForm();

    var $target = $(evt.currentTarget),
        parent = $target.data('in-reply-to');

    if (parent) {
        this.commentForm.detach();
        this.commentForm.appendTo($target.closest('div.comment'));
    }
    else {
        this.commentForm.appendTo(this.$el);
    }
    this.commentForm[0].parent.value = parent || null;
};

Kaiju.prototype.showCommentForm = function() {
    this.commentForm.removeClass('hidden');
};

Kaiju.prototype.hideCommentForm = function() {
    this.commentForm.addClass('hidden');
};

$(function() {
    var kaiju = new Kaiju({
        url: "http://10.4.126.233:2714",
        forum: "5346e494331583002c7de60e",
        page: "local_test_page",
        selector: "section.comments-section"
    });

    kaiju.connect();
});
