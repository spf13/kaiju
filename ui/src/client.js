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
    this.socket.on('commentPosted', _.bind(this.handleCommentPosted, this));

    this._comments = { };

    this.commentTemplate = _.template(
        '<div class="comment" data-id="<%= id %>">' +
        '<div class="comment-header"><strong><%= user.name %></strong> (<%= user.email %>) <%= verb %>:</div>' +
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
    var self = this;
    data = JSON.parse(data);
    if (_.isArray(data)) {
        _.each(data, function(comment) {
            self.renderComment(comment);
        });
    }
};

Kaiju.prototype.handleCommentPosted = function(data) {
    var comment = JSON.parse(data);

    if (comment) {
        this.renderComment(comment);
    }
};

Kaiju.prototype.renderComment = function(comment) {
    var commentTemplate = this.commentTemplate,
        threadTemplate = this.threadTemplate,
        parent = comment.Parent,
        $comment = $(commentTemplate({
            user: {
                name: comment.User.FullName,
                email: comment.User.Email
            },
            verb: comment.Parent ? "replied" : "said",
            body: comment.Body,
            id: comment.Id,
            parentId: comment.Parent || "null"
        }));

    if (parent) {
        var $parentComment = this.$el.find("div.comment[data-id='" + parent + "']"),
            $newThread = $parentComment.find('div.comment-thread');

        if ($newThread.length === 0) {
            $newThread = $('<div class="comment-thread">');
            $parentComment.append($newThread);
        }

        $newThread.append($comment);
    }
    else {
        var $mainThread = this.$el.find('> div.comment-thread');
        if ($mainThread.length === 0) {
            $mainThread = $(threadTemplate());
            this.$el.empty().append($mainThread);
        }
        $mainThread.append($comment);
    }
    $comment.find('a.add-comment').on('click', _.bind(this.onClickShowCommentForm, this));
};

Kaiju.prototype.onSubmitCommentForm = function(evt) {
    evt.stopPropagation();
    evt.preventDefault();

    var form = this.commentForm[0];

    this.postComment({
        fullname: form.user.value,
        email: form.email.value,
        body: form.body.value,
        parent: form.parent.value
    });

    this.commentForm.addClass('hidden');
};

Kaiju.prototype.onClickShowCommentForm = function(evt) {
    evt.stopPropagation();
    evt.preventDefault();

    var $target = $(evt.currentTarget),
        parent = $target.data('in-reply-to');

    if (parent) {
        this.commentForm.detach();
        this.commentForm.insertAfter($target.closest('div.comment').find('> div.comment-actions'));
        this.commentForm[0].parent.value = parent;
    }
    else {
        this.commentForm.appendTo(this.$el);
        this.commentForm[0].parent.value = "";
    }
    this.showCommentForm();
};

Kaiju.prototype.showCommentForm = function() {
    this.commentForm.removeClass('hidden');
    // this.commentForm.find('input, textarea').val('');
};

Kaiju.prototype.hideCommentForm = function() {
    this.commentForm.addClass('hidden');
};
