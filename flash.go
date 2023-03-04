package flash

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Flash struct {
	data   fiber.Map
	config Config
}

type Config struct {
	Name string
}

var DefaultFlash *Flash
var store *session.Store

func init() {
	Default(Config{
		Name: "fiber-app-flash",
	})
}

func Default(config Config) {
	DefaultFlash = New(config)
}

func New(config Config) *Flash {
	store = session.New()
	return &Flash{
		config: config,
		data:   fiber.Map{},
	}
}

func (f *Flash) Get(c *fiber.Ctx) fiber.Map {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	// Get value
	v := sess.Get(f.config.Name)
	vString := ""
	if v != nil {
		vString = v.(string)
	}
	var decoded map[string]interface{}
	json.Unmarshal([]byte(vString), &decoded)

	if decoded != nil {
		f.data = decoded
	} else {
		f.data = fiber.Map{}
	}

	return f.data
}

func (f *Flash) Clear(c *fiber.Ctx) {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Set(f.config.Name, "")

	if err := sess.Save(); err != nil {
		panic(err)
	}
}

func (f *Flash) Redirect(c *fiber.Ctx, location string, data interface{}, status ...int) error {
	f.data = data.(fiber.Map)
	if len(status) > 0 {
		return c.Redirect(location, status[0])
	} else {
		return c.Redirect(location, fiber.StatusFound)
	}
}

func (f *Flash) RedirectToRoute(c *fiber.Ctx, routeName string, data fiber.Map, status ...int) error {
	f.data = data
	if len(status) > 0 {
		return c.RedirectToRoute(routeName, data, status[0])
	} else {
		return c.RedirectToRoute(routeName, data, fiber.StatusFound)
	}
}

func (f *Flash) RedirectBack(c *fiber.Ctx, fallback string, data fiber.Map, status ...int) error {
	f.data = data
	if len(status) > 0 {
		return c.RedirectBack(fallback, status[0])
	} else {
		return c.RedirectBack(fallback, fiber.StatusFound)
	}
}

func (f *Flash) WithError(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.error(c)
	return c
}

func (f *Flash) WithSuccess(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.success(c)
	return c
}

func (f *Flash) WithWarn(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.warn(c)
	return c
}

func (f *Flash) WithInfo(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.info(c)
	return c
}

func (f *Flash) WithData(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.setSession(c)
	return c
}

func (f *Flash) error(c *fiber.Ctx) {
	f.data["error"] = true
	f.setSession(c)
}

func (f *Flash) success(c *fiber.Ctx) {
	f.data["success"] = true
	f.setSession(c)
}

func (f *Flash) warn(c *fiber.Ctx) {
	f.data["warn"] = true
	f.setSession(c)
}

func (f *Flash) info(c *fiber.Ctx) {
	f.data["info"] = true
	f.setSession(c)
}

func (f *Flash) setSession(c *fiber.Ctx) {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	encoded, err := json.Marshal(f.data)
	if err != nil {
		panic(err)
	}

	sess.Set(f.config.Name, string(encoded))

	// Save session
	if err := sess.Save(); err != nil {
		panic(err)
	}
}

func Get(c *fiber.Ctx) fiber.Map {
	return DefaultFlash.Get(c)
}

func Clear(c *fiber.Ctx) {
	DefaultFlash.Clear(c)
}

func Redirect(c *fiber.Ctx, location string, data interface{}, status ...int) error {
	return DefaultFlash.Redirect(c, location, data, status...)
}

func RedirectToRoute(c *fiber.Ctx, routeName string, data fiber.Map, status ...int) error {
	return DefaultFlash.RedirectToRoute(c, routeName, data, status...)
}

func RedirectBack(c *fiber.Ctx, fallback string, data fiber.Map, status ...int) error {
	return DefaultFlash.RedirectBack(c, fallback, data, status...)
}

func WithError(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithError(c, data)
}

func WithSuccess(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithSuccess(c, data)
}

func WithWarn(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithWarn(c, data)
}

func WithInfo(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithInfo(c, data)
}

func WithData(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithData(c, data)
}
