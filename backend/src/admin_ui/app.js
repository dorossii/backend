(function () {
  "use strict";

  var TASK_TAGS = { 0: "掃除", 1: "洗濯", 2: "料理", 4: "その他" };
  var TASK_STATUSES = { 0: "未完了", 1: "承認待ち", 2: "完了" };
  var FRIEND_STATUSES = { 0: "申請中", 1: "成立", 2: "拒否" };

  // ---- 共通ユーティリティ ----

  function el(tag, opts, children) {
    var node = document.createElement(tag);
    opts = opts || {};
    if (opts.className) node.className = opts.className;
    if (opts.text !== undefined) node.textContent = opts.text;
    if (opts.type) node.type = opts.type;
    if (opts.value !== undefined) node.value = opts.value;
    if (opts.placeholder) node.placeholder = opts.placeholder;
    if (opts.checked !== undefined) node.checked = opts.checked;
    if (opts.attrs) {
      Object.keys(opts.attrs).forEach(function (key) {
        node.setAttribute(key, opts.attrs[key]);
      });
    }
    if (opts.onClick) node.addEventListener("click", opts.onClick);
    (children || []).forEach(function (child) {
      if (child) node.appendChild(child);
    });
    return node;
  }

  function clearChildren(node) {
    while (node.firstChild) node.removeChild(node.firstChild);
  }

  var toastEl = document.getElementById("toast");
  var toastTimer = null;

  function showToast(message, isError) {
    clearTimeout(toastTimer);
    toastEl.textContent = message;
    toastEl.hidden = false;
    toastEl.className = "toast " + (isError ? "is-error" : "is-success");
    toastTimer = setTimeout(function () {
      toastEl.hidden = true;
    }, 3500);
  }

  function apiRequest(path, options) {
    options = options || {};
    var fetchOptions = {
      method: options.method || "GET",
      credentials: "same-origin",
      headers: {},
    };
    if (options.body !== undefined) {
      fetchOptions.headers["Content-Type"] = "application/json";
      fetchOptions.body = JSON.stringify(options.body);
    }

    return fetch(path, fetchOptions).then(function (res) {
      if (res.status === 401) {
        window.location.href = "login.html";
        throw new Error("認証切れです");
      }
      if (res.status === 204) return null;
      return res.json().then(function (data) {
        if (!res.ok) {
          throw new Error((data && data.error) || "リクエストに失敗しました");
        }
        return data;
      });
    });
  }

  // ---- モーダル ----

  var modalBackdrop = document.getElementById("modal-backdrop");
  var modalTitle = document.getElementById("modal-title");
  var modalForm = document.getElementById("modal-form");

  function closeModal() {
    modalBackdrop.hidden = true;
    clearChildren(modalForm);
  }

  modalBackdrop.addEventListener("click", function (event) {
    if (event.target === modalBackdrop) closeModal();
  });

  // fields: [{key, label, type: 'text'|'number'|'checkbox'|'select'|'select-string', value, options, hint}]
  function openModal(title, fields, onSubmit) {
    modalTitle.textContent = title;
    clearChildren(modalForm);

    var inputs = {};

    fields.forEach(function (field) {
      if (field.type === "checkbox") {
        var checkboxWrap = el("label", { className: "field checkbox-field" });
        var checkbox = el("input", { type: "checkbox", checked: !!field.value });
        inputs[field.key] = checkbox;
        checkboxWrap.appendChild(checkbox);
        checkboxWrap.appendChild(el("span", { text: field.label }));
        modalForm.appendChild(checkboxWrap);
        return;
      }

      var wrap = el("label", { className: "field" });
      wrap.appendChild(el("span", { className: "field-label", text: field.label }));

      var input;
      if (field.type === "select" || field.type === "select-string") {
        input = document.createElement("select");
        input.className = "input";
        if (field.placeholder) {
          var placeholderOption = el("option", { text: field.placeholder, value: "" });
          placeholderOption.disabled = true;
          if (field.value === undefined || field.value === null || field.value === "") {
            placeholderOption.selected = true;
          }
          input.appendChild(placeholderOption);
        }
        (field.options || []).forEach(function (opt) {
          var optionEl = el("option", { text: opt.label, value: opt.value });
          if (String(opt.value) === String(field.value)) optionEl.selected = true;
          input.appendChild(optionEl);
        });
      } else {
        input = el("input", {
          className: "input",
          type: field.type || "text",
          value: field.value !== undefined && field.value !== null ? field.value : "",
        });
        if (field.readOnly) input.readOnly = true;
      }
      inputs[field.key] = input;
      wrap.appendChild(input);
      if (field.hint) wrap.appendChild(el("span", { className: "field-hint", text: field.hint }));
      modalForm.appendChild(wrap);
    });

    var actions = el("div", { className: "modal-actions" }, [
      el("button", {
        className: "btn btn-secondary",
        type: "button",
        text: "キャンセル",
        onClick: closeModal,
      }),
      el("button", { className: "btn btn-primary", type: "submit", text: "保存" }),
    ]);
    modalForm.appendChild(actions);

    modalForm.onsubmit = function (event) {
      event.preventDefault();
      var values = {};
      fields.forEach(function (field) {
        var input = inputs[field.key];
        if (field.type === "checkbox") {
          values[field.key] = input.checked;
        } else if (field.type === "number") {
          values[field.key] = input.value === "" ? null : Number(input.value);
        } else if (field.type === "select") {
          values[field.key] = Number(input.value);
        } else {
          values[field.key] = input.value;
        }
      });
      onSubmit(values);
    };

    modalBackdrop.hidden = false;
  }

  // ---- ナビゲーション ----

  var navItems = document.querySelectorAll(".nav-item");
  var views = document.querySelectorAll(".view");
  var viewTitle = document.getElementById("view-title");
  var viewLabels = {
    users: "ユーザー",
    "base-tasks": "タスクテンプレート",
    tasks: "タスク",
    friendships: "フレンド関係",
  };
  var viewLoaders = {
    users: loadUsers,
    "base-tasks": loadBaseTasks,
    tasks: loadTasks,
    friendships: renderFriendships,
  };

  navItems.forEach(function (item) {
    item.addEventListener("click", function () {
      var target = item.getAttribute("data-view");

      navItems.forEach(function (n) {
        n.classList.toggle("is-active", n === item);
      });
      views.forEach(function (view) {
        view.classList.toggle("is-active", view.id === "view-" + target);
      });
      viewTitle.textContent = viewLabels[target];
      viewLoaders[target]();
    });
  });

  document.getElementById("logout-btn").addEventListener("click", function () {
    apiRequest("logout", { method: "POST" }).then(function () {
      window.location.href = "login.html";
    });
  });

  // ---- ユーザー管理 ----

  var usersTbody = document.getElementById("users-tbody");

  function loadUsers() {
    var search = document.getElementById("user-search").value.trim();
    var query = search ? "?search=" + encodeURIComponent(search) : "";

    apiRequest("users" + query)
      .then(function (data) {
        renderUsers((data && data.users) || []);
      })
      .catch(function (err) {
        showToast(err.message, true);
      });
  }

  function renderUsers(users) {
    clearChildren(usersTbody);

    if (users.length === 0) {
      usersTbody.appendChild(
        el("tr", { className: "empty-row" }, [el("td", { attrs: { colspan: "7" }, text: "ユーザーが見つかりません" })])
      );
      return;
    }

    users.forEach(function (user) {
      var editBtn = el("button", {
        className: "btn btn-secondary btn-sm",
        text: "HP/汚れを編集",
        onClick: function () {
          openEditUserStatsModal(user);
        },
      });

      var row = el("tr", {}, [
        el("td", { text: user.UserID }),
        el("td", { text: user.UserName }),
        el("td", { text: user.Mailadress }),
        el("td", { text: String(user.HealthPoint) }),
        el("td", { text: String(user.DirtLevel) }),
        el("td", { text: String(user.Combo) }),
        el("td", { className: "row-actions" }, [editBtn]),
      ]);
      usersTbody.appendChild(row);
    });
  }

  function openEditUserStatsModal(user) {
    openModal(
      "HP / 汚れレベルの編集: " + user.UserName,
      [
        { key: "HealthPoint", label: "HealthPoint", type: "number", value: user.HealthPoint },
        { key: "DirtLevel", label: "DirtLevel", type: "number", value: user.DirtLevel },
      ],
      function (values) {
        apiRequest("users/" + encodeURIComponent(user.UserID) + "/stats", {
          method: "PUT",
          body: values,
        })
          .then(function () {
            closeModal();
            showToast("更新しました", false);
            loadUsers();
          })
          .catch(function (err) {
            showToast(err.message, true);
          });
      }
    );
  }

  document.getElementById("user-search-btn").addEventListener("click", loadUsers);
  document.getElementById("user-search").addEventListener("keydown", function (event) {
    if (event.key === "Enter") loadUsers();
  });

  // ---- タスクテンプレート(BaseTask) ----

  var baseTasksTbody = document.getElementById("base-tasks-tbody");

  function loadBaseTasks() {
    apiRequest("base-tasks")
      .then(function (data) {
        renderBaseTasks((data && data.baseTasks) || []);
      })
      .catch(function (err) {
        showToast(err.message, true);
      });
  }

  function renderBaseTasks(baseTasks) {
    clearChildren(baseTasksTbody);

    if (baseTasks.length === 0) {
      baseTasksTbody.appendChild(
        el("tr", { className: "empty-row" }, [el("td", { attrs: { colspan: "8" }, text: "テンプレートがありません" })])
      );
      return;
    }

    baseTasks.forEach(function (bt) {
      var editBtn = el("button", {
        className: "btn btn-secondary btn-sm",
        text: "編集",
        onClick: function () {
          openBaseTaskModal(bt);
        },
      });
      var deleteBtn = el("button", {
        className: "btn btn-danger btn-sm",
        text: "削除",
        onClick: function () {
          if (!window.confirm("「" + bt.TaskName + "」を削除しますか？")) return;
          apiRequest("base-tasks/" + encodeURIComponent(bt.BaseID), { method: "DELETE" })
            .then(function () {
              showToast("削除しました", false);
              loadBaseTasks();
            })
            .catch(function (err) {
              showToast(err.message, true);
            });
        },
      });

      var row = el("tr", {}, [
        el("td", { text: bt.BaseID }),
        el("td", { text: bt.TaskName }),
        el("td", { text: bt.Description }),
        el("td", { text: String(bt.DifficultyLevel) }),
        el("td", { text: String(bt.DueTime) }),
        el("td", { text: bt.ImageFlag ? "必須" : "任意" }),
        el("td", { text: TASK_TAGS[bt.Tags] || String(bt.Tags) }),
        el("td", { className: "row-actions" }, [editBtn, deleteBtn]),
      ]);
      baseTasksTbody.appendChild(row);
    });
  }

  function tagOptions() {
    return Object.keys(TASK_TAGS).map(function (key) {
      return { label: TASK_TAGS[key], value: Number(key) };
    });
  }

  function openBaseTaskModal(baseTask) {
    var isNew = !baseTask;
    var fields = [
      { key: "BaseID", label: "BaseID", type: "text", value: baseTask ? baseTask.BaseID : "", readOnly: !isNew },
      { key: "TaskName", label: "タスク名", type: "text", value: baseTask ? baseTask.TaskName : "" },
      { key: "Description", label: "説明", type: "text", value: baseTask ? baseTask.Description : "" },
      { key: "DifficultyLevel", label: "難易度", type: "number", value: baseTask ? baseTask.DifficultyLevel : 1 },
      { key: "DueTime", label: "期限(日)", type: "number", value: baseTask ? baseTask.DueTime : 1 },
      { key: "ImageFlag", label: "画像必須", type: "checkbox", value: baseTask ? baseTask.ImageFlag : false },
      { key: "Tags", label: "タグ", type: "select", value: baseTask ? baseTask.Tags : 0, options: tagOptions() },
    ];

    openModal(isNew ? "タスクテンプレートの新規作成" : "タスクテンプレートの編集", fields, function (values) {
      var path = isNew
        ? "base-tasks"
        : "base-tasks/" + encodeURIComponent(baseTask.BaseID);
      var method = isNew ? "POST" : "PUT";

      apiRequest(path, { method: method, body: values })
        .then(function () {
          closeModal();
          showToast("保存しました", false);
          loadBaseTasks();
        })
        .catch(function (err) {
          showToast(err.message, true);
        });
    });
  }

  document.getElementById("base-task-new-btn").addEventListener("click", function () {
    openBaseTaskModal(null);
  });

  // ---- タスク(Task) ----

  var tasksTbody = document.getElementById("tasks-tbody");

  function loadTasks() {
    var userID = document.getElementById("task-user-filter").value.trim();
    var query = userID ? "?userID=" + encodeURIComponent(userID) : "";

    apiRequest("tasks" + query)
      .then(function (data) {
        renderTasks((data && data.tasks) || []);
      })
      .catch(function (err) {
        showToast(err.message, true);
      });
  }

  function statusBadgeClass(status) {
    if (status === 1) return "badge badge-pending";
    if (status === 2) return "badge badge-completed";
    return "badge badge-incomplete";
  }

  function renderTasks(tasks) {
    clearChildren(tasksTbody);

    if (tasks.length === 0) {
      tasksTbody.appendChild(
        el("tr", { className: "empty-row" }, [el("td", { attrs: { colspan: "8" }, text: "タスクが見つかりません" })])
      );
      return;
    }

    tasks.forEach(function (task) {
      var statusBadge = el("span", {
        className: statusBadgeClass(task.Status),
        text: TASK_STATUSES[task.Status] || String(task.Status),
      });

      var imageLink = null;
      if (task.ImageID) {
        imageLink = el("a", {
          text: "確認",
          attrs: { href: "tasks/" + encodeURIComponent(task.TaskID) + "/image", target: "_blank", rel: "noopener" },
        });
      } else {
        imageLink = el("span", { text: "-" });
      }

      var editBtn = el("button", {
        className: "btn btn-secondary btn-sm",
        text: "編集",
        onClick: function () {
          openTaskEditModal(task);
        },
      });
      var statusBtn = el("button", {
        className: "btn btn-secondary btn-sm",
        text: "状態変更",
        onClick: function () {
          openTaskStatusModal(task);
        },
      });
      var deleteBtn = el("button", {
        className: "btn btn-danger btn-sm",
        text: "削除",
        onClick: function () {
          if (!window.confirm("このタスクを削除しますか？")) return;
          apiRequest("tasks/" + encodeURIComponent(task.TaskID), { method: "DELETE" })
            .then(function () {
              showToast("削除しました", false);
              loadTasks();
            })
            .catch(function (err) {
              showToast(err.message, true);
            });
        },
      });

      var row = el("tr", {}, [
        el("td", { text: task.TaskID }),
        el("td", { text: task.UserID }),
        el("td", { text: task.BaseID }),
        el("td", {}, [statusBadge]),
        el("td", { text: task.RequireImage ? "必須" : "任意" }),
        el("td", {}, [imageLink]),
        el("td", { text: task.Message || "" }),
        el("td", { className: "row-actions" }, [editBtn, statusBtn, deleteBtn]),
      ]);
      tasksTbody.appendChild(row);
    });
  }

  // ユーザー一覧を取得し、セレクトボックス用のoptions配列に変換する
  function fetchUserOptions() {
    return apiRequest("users").then(function (data) {
      var users = (data && data.users) || [];
      return users.map(function (user) {
        return { label: user.UserName + "（" + user.UserID + "）", value: user.UserID };
      });
    });
  }

  function openTaskNewModal() {
    Promise.all([fetchUserOptions(), apiRequest("base-tasks")])
      .then(function (results) {
        var userOptions = results[0] || [];
        var baseTasks = (results[1] && results[1].baseTasks) || [];

        var baseTaskOptions = baseTasks.map(function (bt) {
          return { label: bt.TaskName + "（" + TASK_TAGS[bt.Tags] + " / 難易度" + bt.DifficultyLevel + "）", value: bt.BaseID };
        });

        openModal(
          "タスクの手動生成",
          [
            {
              key: "UserID",
              label: "対象ユーザー",
              type: "select-string",
              value: "",
              options: userOptions,
              placeholder: userOptions.length ? "ユーザーを選択してください" : "ユーザーが登録されていません",
            },
            {
              key: "BaseID",
              label: "テンプレート",
              type: "select-string",
              value: "",
              options: baseTaskOptions,
              placeholder: baseTaskOptions.length ? "テンプレートを選択してください" : "テンプレートが登録されていません",
            },
            { key: "RequireImage", label: "画像を必須にする", type: "checkbox", value: false },
            { key: "Message", label: "メッセージ(任意)", type: "text", value: "" },
          ],
          function (values) {
            if (!values.UserID || !values.BaseID) {
              showToast("対象ユーザーとテンプレートを選択してください", true);
              return;
            }
            apiRequest("tasks", { method: "POST", body: values })
              .then(function () {
                closeModal();
                showToast("タスクを作成しました", false);
                loadTasks();
              })
              .catch(function (err) {
                showToast(err.message, true);
              });
          }
        );
      })
      .catch(function (err) {
        showToast(err.message, true);
      });
  }

  function openTaskEditModal(task) {
    openModal(
      "タスクの編集",
      [
        { key: "RequireImage", label: "画像必須", type: "checkbox", value: task.RequireImage },
        { key: "Message", label: "メッセージ", type: "text", value: task.Message || "" },
      ],
      function (values) {
        var body = Object.assign({}, task, values);
        apiRequest("tasks/" + encodeURIComponent(task.TaskID), { method: "PUT", body: body })
          .then(function () {
            closeModal();
            showToast("更新しました", false);
            loadTasks();
          })
          .catch(function (err) {
            showToast(err.message, true);
          });
      }
    );
  }

  function openTaskStatusModal(task) {
    var options = Object.keys(TASK_STATUSES).map(function (key) {
      return { label: TASK_STATUSES[key], value: Number(key) };
    });

    openModal(
      "タスクの状態変更",
      [{ key: "Status", label: "状態", type: "select", value: task.Status, options: options }],
      function (values) {
        apiRequest("tasks/" + encodeURIComponent(task.TaskID) + "/status", {
          method: "PUT",
          body: values,
        })
          .then(function () {
            closeModal();
            showToast("状態を更新しました", false);
            loadTasks();
          })
          .catch(function (err) {
            showToast(err.message, true);
          });
      }
    );
  }

  document.getElementById("task-new-btn").addEventListener("click", openTaskNewModal);
  document.getElementById("task-filter-btn").addEventListener("click", loadTasks);
  document.getElementById("task-user-filter").addEventListener("keydown", function (event) {
    if (event.key === "Enter") loadTasks();
  });

  // ---- フレンド関係 ----

  var friendshipsTbody = document.getElementById("friendships-tbody");
  var friendshipRows = [];

  function renderFriendships() {
    clearChildren(friendshipsTbody);

    if (friendshipRows.length === 0) {
      friendshipsTbody.appendChild(
        el("tr", { className: "empty-row" }, [
          el("td", { attrs: { colspan: "4" }, text: "「フレンド関係を作成」から操作対象を追加してください" }),
        ])
      );
      return;
    }

    friendshipRows.forEach(function (pair) {
      var statusBtn = el("button", {
        className: "btn btn-secondary btn-sm",
        text: "状態変更",
        onClick: function () {
          openFriendshipStatusModal(pair);
        },
      });
      var deleteBtn = el("button", {
        className: "btn btn-danger btn-sm",
        text: "削除",
        onClick: function () {
          if (!window.confirm("このフレンド関係を削除しますか？")) return;
          apiRequest(
            "friendships/" + encodeURIComponent(pair.userID) + "/" + encodeURIComponent(pair.friendID),
            { method: "DELETE" }
          )
            .then(function () {
              friendshipRows = friendshipRows.filter(function (row) {
                return row !== pair;
              });
              showToast("削除しました", false);
              renderFriendships();
            })
            .catch(function (err) {
              showToast(err.message, true);
            });
        },
      });

      var row = el("tr", {}, [
        el("td", { text: pair.userID }),
        el("td", { text: pair.friendID }),
        el("td", { text: FRIEND_STATUSES[pair.status] || "-" }),
        el("td", { className: "row-actions" }, [statusBtn, deleteBtn]),
      ]);
      friendshipsTbody.appendChild(row);
    });
  }

  function openFriendshipNewModal() {
    fetchUserOptions()
      .then(function (userOptions) {
        openModal(
          "フレンド関係の作成(強制成立)",
          [
            {
              key: "UserID",
              label: "ユーザー",
              type: "select-string",
              value: "",
              options: userOptions,
              placeholder: userOptions.length ? "ユーザーを選択してください" : "ユーザーが登録されていません",
            },
            {
              key: "FriendID",
              label: "フレンド",
              type: "select-string",
              value: "",
              options: userOptions,
              placeholder: userOptions.length ? "フレンドを選択してください" : "ユーザーが登録されていません",
            },
          ],
          function (values) {
            if (!values.UserID || !values.FriendID) {
              showToast("ユーザーとフレンドを選択してください", true);
              return;
            }
            if (values.UserID === values.FriendID) {
              showToast("ユーザーとフレンドには異なるユーザーを選択してください", true);
              return;
            }
            openFriendshipCreateRequest(values);
          }
        );
      })
      .catch(function (err) {
        showToast(err.message, true);
      });
  }

  function openFriendshipCreateRequest(values) {
    apiRequest("friendships", { method: "POST", body: values })
      .then(function () {
        closeModal();
        showToast("フレンド関係を成立させました", false);
        friendshipRows.push({ userID: values.UserID, friendID: values.FriendID, status: 1 });
        renderFriendships();
      })
      .catch(function (err) {
        showToast(err.message, true);
      });
  }

  function openFriendshipStatusModal(pair) {
    var options = Object.keys(FRIEND_STATUSES).map(function (key) {
      return { label: FRIEND_STATUSES[key], value: Number(key) };
    });

    openModal(
      "フレンド関係の状態変更",
      [{ key: "Status", label: "状態", type: "select", value: pair.status, options: options }],
      function (values) {
        apiRequest(
          "friendships/" + encodeURIComponent(pair.userID) + "/" + encodeURIComponent(pair.friendID),
          { method: "PUT", body: values }
        )
          .then(function () {
            pair.status = values.Status;
            closeModal();
            showToast("状態を更新しました", false);
            renderFriendships();
          })
          .catch(function (err) {
            showToast(err.message, true);
          });
      }
    );
  }

  document.getElementById("friendship-new-btn").addEventListener("click", openFriendshipNewModal);

  // ---- 初期表示 ----
  loadUsers();
})();
