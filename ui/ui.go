package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/leinonen/bbs/domain"
	"github.com/leinonen/bbs/repository"
	"golang.org/x/term"
)

type UI struct {
	term    *term.Terminal
	repos   *repository.Manager
	session *domain.Session
}

func NewUI(term *term.Terminal, repos *repository.Manager, session *domain.Session) *UI {
	return &UI{
		term:    term,
		repos:   repos,
		session: session,
	}
}

func (ui *UI) Run() {
	ui.clear()
	ui.showWelcome()

	for {
		if ui.session.User == nil || ui.session.User.ID == 0 {
			if !ui.showLoginMenu() {
				return
			}
		} else {
			if !ui.showMainMenu() {
				return
			}
		}
	}
}

func (ui *UI) showWelcome() {
	ui.printHeader("Welcome to Go BBS System")
	ui.println("")
	ui.println("A modern take on the classic Bulletin Board System")
	ui.println("Connected via SSH")
	ui.println("")
	ui.printLine()
}

func (ui *UI) showLoginMenu() bool {
	ui.println("")
	ui.printHeader("Login Menu")
	ui.println("1. Login")
	ui.println("2. Register")
	ui.println("3. Continue as Guest")
	ui.println("4. Exit")
	ui.println("")

	choice := ui.readLine("Select option: ")

	switch choice {
	case "1":
		ui.handleLogin()
	case "2":
		ui.handleRegister()
	case "3":
		ui.session.User = &domain.User{
			ID:       0,
			Username: "guest",
		}
	case "4":
		ui.println("Goodbye!")
		return false
	default:
		ui.printError("Invalid option")
	}

	return true
}

func (ui *UI) showMainMenu() bool {
	ui.clear()
	ui.printHeader(fmt.Sprintf("Main Menu - Welcome %s", ui.session.User.Username))
	ui.println("")
	ui.println("1. Browse Boards")
	ui.println("2. Recent Posts")
	ui.println("3. Search")
	ui.println("4. User Profile")
	ui.println("5. Who's Online")
	if ui.session.User.IsAdmin {
		ui.println("6. Admin Panel")
	}
	ui.println("9. Logout")
	ui.println("0. Exit")
	ui.println("")

	choice := ui.readLine("Select option: ")

	switch choice {
	case "1":
		ui.browseBoards()
	case "2":
		ui.showRecentPosts()
	case "3":
		ui.search()
	case "4":
		ui.showProfile()
	case "5":
		ui.showOnlineUsers()
	case "6":
		if ui.session.User.IsAdmin {
			ui.adminPanel()
		}
	case "9":
		ui.session.User = nil
		ui.println("Logged out successfully")
		time.Sleep(1 * time.Second)
	case "0":
		ui.println("Goodbye!")
		return false
	default:
		ui.printError("Invalid option")
	}

	return true
}

func (ui *UI) handleLogin() {
	ui.clear()
	ui.printHeader("Login")

	username := ui.readLine("Username: ")
	ui.term.SetPrompt("Password: ")
	password, _ := ui.term.ReadPassword("Password: ")
	ui.term.SetPrompt("")

	user, err := ui.repos.User.Authenticate(username, password)
	if err != nil {
		ui.printError("Invalid credentials")
		time.Sleep(2 * time.Second)
		return
	}

	ui.repos.User.UpdateLastLogin(user.ID)
	ui.session.User = user
	ui.printSuccess(fmt.Sprintf("Welcome back, %s!", user.Username))
	time.Sleep(1 * time.Second)
}

func (ui *UI) handleRegister() {
	ui.clear()
	ui.printHeader("Register New Account")

	username := ui.readLine("Username: ")
	email := ui.readLine("Email: ")

	ui.term.SetPrompt("Password: ")
	password, _ := ui.term.ReadPassword("Password: ")
	ui.term.SetPrompt("Confirm Password: ")
	confirm, _ := ui.term.ReadPassword("Confirm Password: ")
	ui.term.SetPrompt("")

	if password != confirm {
		ui.printError("Passwords do not match")
		time.Sleep(2 * time.Second)
		return
	}

	user := domain.NewUser(username, email)
	user.Password = password

	err := ui.repos.User.Create(user)
	if err != nil {
		ui.printError(fmt.Sprintf("Registration failed: %v", err))
		time.Sleep(2 * time.Second)
		return
	}

	ui.session.User = user
	ui.printSuccess("Registration successful!")
	time.Sleep(1 * time.Second)
}

func (ui *UI) browseBoards() {
	for {
		ui.clear()
		ui.printHeader("Message Boards")

		boards, err := ui.repos.Board.GetAll()
		if err != nil {
			ui.printError(fmt.Sprintf("Error loading boards: %v", err))
			return
		}

		for i, board := range boards {
			ui.println(fmt.Sprintf("%d. [%s] %s (%d posts)",
				i+1, board.Name, board.Description, board.PostCount))
		}

		ui.println("")
		ui.println("Enter board number (0 to go back): ")

		choice := ui.readLine("> ")
		if choice == "0" {
			return
		}

		num, err := strconv.Atoi(choice)
		if err != nil || num < 1 || num > len(boards) {
			ui.printError("Invalid selection")
			continue
		}

		ui.viewBoard(boards[num-1])
	}
}

func (ui *UI) viewBoard(board *domain.Board) {
	page := 0
	pageSize := 20

	for {
		ui.clear()
		ui.printHeader(fmt.Sprintf("Board: %s", board.Name))
		ui.println(board.Description)
		ui.printLine()

		posts, err := ui.repos.Post.GetByBoard(board.ID, pageSize, page*pageSize)
		if err != nil {
			ui.printError(fmt.Sprintf("Error loading posts: %v", err))
			return
		}

		if len(posts) == 0 {
			ui.println("No posts yet. Be the first to post!")
		} else {
			for i, post := range posts {
				ui.println(fmt.Sprintf("%d. %s - by %s (%d replies)",
					i+1, post.Title, post.Username, post.Replies))
				ui.println(fmt.Sprintf("   %s", ui.formatTime(post.CreatedAt)))
			}
		}

		ui.println("")
		ui.println("Commands: (N)ew post, (V)iew post #, (R)efresh, (B)ack")
		if page > 0 {
			ui.print(", (P)revious page")
		}
		if len(posts) == pageSize {
			ui.print(", (F)orward page")
		}
		ui.println("")

		cmd := ui.readLine("> ")
		cmd = strings.ToLower(strings.TrimSpace(cmd))

		switch {
		case cmd == "b":
			return
		case cmd == "n":
			if ui.session.User == nil || ui.session.User.ID == 0 {
				ui.printError("Please login to post")
				time.Sleep(2 * time.Second)
			} else {
				ui.createPost(board.ID, nil)
			}
		case cmd == "r":
			continue
		case cmd == "p" && page > 0:
			page--
		case cmd == "f" && len(posts) == pageSize:
			page++
		case strings.HasPrefix(cmd, "v"):
			parts := strings.Fields(cmd)
			if len(parts) == 2 {
				if num, err := strconv.Atoi(parts[1]); err == nil && num > 0 && num <= len(posts) {
					ui.viewPost(posts[num-1])
				}
			}
		default:
			if num, err := strconv.Atoi(cmd); err == nil && num > 0 && num <= len(posts) {
				ui.viewPost(posts[num-1])
			}
		}
	}
}

func (ui *UI) viewPost(post *domain.Post) {
	ui.clear()
	ui.printHeader(post.Title)
	ui.println(fmt.Sprintf("Posted by %s on %s", post.Username, ui.formatTime(post.CreatedAt)))
	ui.printLine()
	ui.println(post.Content)
	ui.printLine()

	replies, _ := ui.repos.Post.GetReplies(post.ID)
	if len(replies) > 0 {
		ui.println(fmt.Sprintf("--- %d Replies ---", len(replies)))
		for _, reply := range replies {
			ui.println("")
			ui.println(fmt.Sprintf("By %s on %s:", reply.Username, ui.formatTime(reply.CreatedAt)))
			ui.println(reply.Content)
		}
		ui.printLine()
	}

	ui.println("")
	ui.println("Commands: (R)eply, (B)ack")

	cmd := ui.readLine("> ")
	cmd = strings.ToLower(strings.TrimSpace(cmd))

	if cmd == "r" {
		if ui.session.User == nil || ui.session.User.ID == 0 {
			ui.printError("Please login to reply")
			time.Sleep(2 * time.Second)
		} else {
			ui.createPost(post.BoardID, &post.ID)
		}
	}
}

func (ui *UI) createPost(boardID int, replyTo *int) {
	ui.clear()
	if replyTo != nil {
		ui.printHeader("Write Reply")
	} else {
		ui.printHeader("New Post")
	}

	var title string
	if replyTo == nil {
		title = ui.readLine("Title: ")
	}

	ui.println("Content (type '.' on a new line to finish):")
	lines := []string{}
	for {
		line := ui.readLine("")
		if line == "." {
			break
		}
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")

	if title == "" && replyTo == nil {
		ui.printError("Title cannot be empty")
		return
	}

	if content == "" {
		ui.printError("Content cannot be empty")
		return
	}

	var post *domain.Post
	if replyTo != nil {
		post = domain.NewReply(boardID, ui.session.User.ID, ui.session.User.Username, content, *replyTo)
	} else {
		post = domain.NewPost(boardID, ui.session.User.ID, ui.session.User.Username, title, content)
	}

	err := ui.repos.Post.Create(post)
	if err != nil {
		ui.printError(fmt.Sprintf("Failed to create post: %v", err))
	} else {
		ui.printSuccess("Post created successfully!")
	}
	time.Sleep(1 * time.Second)
}

func (ui *UI) showRecentPosts() {
	ui.clear()
	ui.printHeader("Recent Posts")

	posts, err := ui.repos.Post.GetRecent(20)
	if err != nil {
		ui.printError(fmt.Sprintf("Error loading posts: %v", err))
		return
	}

	for i, post := range posts {
		ui.println(fmt.Sprintf("%d. %s - by %s", i+1, post.Title, post.Username))
		ui.println(fmt.Sprintf("   %s", ui.formatTime(post.CreatedAt)))
	}

	ui.println("")
	ui.readLine("Press Enter to continue...")
}

func (ui *UI) search() {
	ui.clear()
	ui.printHeader("Search")
	ui.println("Search functionality coming soon...")
	ui.readLine("Press Enter to continue...")
}

func (ui *UI) showProfile() {
	ui.clear()
	ui.printHeader("User Profile")
	ui.println(fmt.Sprintf("Username: %s", ui.session.User.Username))
	ui.println(fmt.Sprintf("Email: %s", ui.session.User.Email))
	ui.println(fmt.Sprintf("Member since: %s", ui.formatTime(ui.session.User.CreatedAt)))
	ui.println(fmt.Sprintf("Last login: %s", ui.formatTime(ui.session.User.LastLogin)))
	if ui.session.User.IsAdmin {
		ui.println("Status: Administrator")
	}
	ui.println("")
	ui.readLine("Press Enter to continue...")
}

func (ui *UI) showOnlineUsers() {
	ui.clear()
	ui.printHeader("Online Users")
	ui.println("Online users functionality coming soon...")
	ui.readLine("Press Enter to continue...")
}

func (ui *UI) adminPanel() {
	ui.clear()
	ui.printHeader("Admin Panel")
	ui.println("1. Create Board")
	ui.println("2. Manage Users")
	ui.println("3. System Stats")
	ui.println("0. Back")

	choice := ui.readLine("Select option: ")

	switch choice {
	case "1":
		ui.createBoard()
	case "2":
		ui.println("User management coming soon...")
		ui.readLine("Press Enter to continue...")
	case "3":
		ui.println("System stats coming soon...")
		ui.readLine("Press Enter to continue...")
	}
}

func (ui *UI) createBoard() {
	ui.clear()
	ui.printHeader("Create New Board")

	name := ui.readLine("Board name: ")
	description := ui.readLine("Description: ")

	board := domain.NewBoard(name, description)
	err := ui.repos.Board.Create(board)
	if err != nil {
		ui.printError(fmt.Sprintf("Failed to create board: %v", err))
	} else {
		ui.printSuccess("Board created successfully!")
	}
	time.Sleep(2 * time.Second)
}
