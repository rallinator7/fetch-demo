tilt_dir = os.getcwd()

local_resource(
  'mage:unit',
  cmd = ['mage', '-v', 'unit'],
  dir=os.path.join(tilt_dir),
  auto_init = False,
  trigger_mode = TRIGGER_MODE_MANUAL,
  labels=['mage']
)

local_resource(
  'mage:build',
  cmd = ['mage', '-v', 'build'],
  dir=os.path.join(tilt_dir),
  deps=[
          os.path.join(tilt_dir, 'cmd', 'payer', 'main.go'),
          os.path.join(tilt_dir, 'cmd', 'user', 'main.go'),
          os.path.join(tilt_dir,'internal'),
      ],
  labels=['mage']
)

apps = ['payer', 'points', 'user']
for app in apps:
  docker_build('%s-service:latest' % app,
              context=os.path.join(tilt_dir, 'cmd', app),
              dockerfile=os.path.join(tilt_dir, 'cmd', app ,'Dockerfile'),
              only=[
                  os.path.join(tilt_dir,'cmd', app),
                  os.path.join(tilt_dir,'internal', app),
                  os.path.join(tilt_dir,'internal', 'environment'),
                  os.path.join(tilt_dir,'internal', 'logger'),
                  os.path.join(tilt_dir, 'go.mod'),
                  os.path.join(tilt_dir, 'go.sum'),
              ],
  )


docker_compose(os.path.join(tilt_dir, 'docker-compose.apps.yaml'))